package elecard

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"io"
	"regexp"
)

type Service interface {
	CreateTask(ctx context.Context, task CreateTaskRequest) (CreateTaskResponse, error)
	GetStatus(ctx context.Context, status GetStatusRequest, fileName string) (string, error)
}

type service struct {
	client *websocket.Conn
	lg     *zap.SugaredLogger
}

func NewService(client *websocket.Conn, lg *zap.SugaredLogger) *service {
	return &service{client: client, lg: lg}
}

func (s *service) CreateTask(ctx context.Context, task CreateTaskRequest) (CreateTaskResponse, error) {
	taskRequest, err := xml.Marshal(task)
	if err != nil {
		return CreateTaskResponse{}, err
	}

	if err := s.client.WriteMessage(websocket.TextMessage, append(XmlHeader, taskRequest...)); err != nil {
		return CreateTaskResponse{}, err
	}

	for {
		var taskResponse CreateTaskResponse
		_, msg, err := s.client.ReadMessage()
		if err != nil {
			return CreateTaskResponse{}, err
		}

		if err := s.decodeXml(msg, &taskResponse); err != nil {
			fmt.Println("unmarshal:", err)
			return CreateTaskResponse{}, err
		}
		return taskResponse, nil
	}
}

func (s *service) GetStatus(ctx context.Context, req GetStatusRequest, fileName string) (string, error) {
	statusRequest, err := xml.Marshal(req)
	if err != nil {
		return "", err
	}

	for {
		if err := s.client.WriteMessage(websocket.TextMessage, append(XmlHeader, statusRequest...)); err != nil {
			return "", err
		}

		var statusCode string
		var statusResponse GetStatusResponse
		_, msg, err := s.client.ReadMessage()
		if err != nil {
			return "", err
		}

		if err := s.decodeXml(msg, &statusResponse); err != nil {
			fmt.Println("unmarshal:", err)
			return "", err
		}
		reg := regexp.MustCompile(`\[(.*?)\]` + " " + fileName)
		matches := reg.FindAllStringSubmatch(statusResponse.GetValue.RetVal, -1)
		for i := range matches {
			statusCode = matches[i][1]
		}

		s.lg.Infof("file[%s] status is %s", fileName, statusCode)
		switch statusCode {
		case SuccessStatus:
			return statusCode, nil
		case FaultStatus:
			return statusCode, errors.New("status is Fault")
		case CriticalErrorStatus:
			return statusCode, errors.New("status is Critical Error")
		}
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
