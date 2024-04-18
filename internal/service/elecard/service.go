package elecard

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/jashakimov/converter/internal/service/socket"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"io"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	CreateTask(ctx context.Context, task CreateTaskRequest) (CreateTaskResponse, error)
	GetStatus(ctx context.Context, status GetStatusRequest, fileName string, delaySec time.Duration) (string, error)
}

type service struct {
	client socket.Service
	lg     *zap.SugaredLogger
}

func NewService(client socket.Service, lg *zap.SugaredLogger) Service {
	return &service{client: client, lg: lg}
}

func (s *service) CreateTask(ctx context.Context, task CreateTaskRequest) (CreateTaskResponse, error) {
	taskRequest, err := xml.Marshal(task)
	if err != nil {
		return CreateTaskResponse{}, err
	}

	s.lg.Info("Добавляем файл в задачу в Elecard:", string(append(XmlHeader, taskRequest...)))
	response, err := s.client.Send(append(XmlHeader, taskRequest...))
	if err != nil {
		s.lg.Error("Ошибка из Elecard:", err.Error())
		return CreateTaskResponse{}, err
	}

	var taskResponse CreateTaskResponse
	if err := xml.Unmarshal(response, &taskResponse); err != nil {
		s.lg.Error("Ошибка парсинга:", err.Error(), string(response))
		return CreateTaskResponse{}, err
	}

	s.lg.Info("Получили ответ от Elecard:", string(response))
	return taskResponse, nil

}

func (s *service) GetStatus(ctx context.Context, req GetStatusRequest, requestPath string, delaySec time.Duration) (string, error) {
	splitted := strings.Split(requestPath, "\\")
	fileName := splitted[len(splitted)-1]

	statusRequest, err := xml.Marshal(req)
	if err != nil {
		return "", err
	}

	for {
		s.lg.Info("Проверяем статус в Elecard:", string(append(XmlHeader, statusRequest...)))
		response, err := s.client.Send(append(XmlHeader, statusRequest...))
		if err != nil {
			s.lg.Error("Ошибка из Elecard:", err.Error())
			return "", err
		}

		var statusCode string
		var statusResponse GetStatusResponse
		if err := xml.Unmarshal(response, &statusResponse); err != nil {
			s.lg.Error("Ошибка парсинга: ", err.Error(), string(response))
			return "", err
		}
		s.lg.Info("Получили ответ от Elecard:", string(response))
		reg := regexp.MustCompile(`\[(.*?)\]` + " " + fileName)
		matches := reg.FindAllStringSubmatch(statusResponse.GetValue.RetVal, -1)
		for i := range matches {
			statusCode = matches[i][1]
		}

		s.lg.Infof("Статус для файла[%s] - [%s]", fileName, statusCode)
		switch statusCode {
		case SuccessStatus:
			return statusCode, nil
		case FaultStatus:
			return statusCode, errors.New("status is Fault")
		case CriticalErrorStatus:
			return statusCode, errors.New("status is Critical Error")
		case NotFoundStatus:
			return statusCode, errors.New("file not found")
		}
		s.lg.Infof("Повторям запрос через %d сек", delaySec)
		time.Sleep(delaySec)
	}
}

func (s *service) decodeXml(b []byte, obj any) error {
	d := xml.NewDecoder(bytes.NewReader(b))
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unhandled charset: %s", charset)
	}
	return d.Decode(obj)
}
