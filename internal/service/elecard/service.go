package elecard

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"io"
)

type Service interface {
	CreateTask(ctx context.Context, task CreateTaskRequest) (CreateTaskResponse, error)
	GetStatus(ctx context.Context, id string)
}

type service struct {
	client *websocket.Conn
	lg     *zap.SugaredLogger
}

func NewService(client *websocket.Conn, lg *zap.SugaredLogger) Service {
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

		d := xml.NewDecoder(bytes.NewReader(msg))
		d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			if charset == "windows-1251" {
				return charmap.Windows1251.NewDecoder().Reader(input), nil
			}
			return nil, fmt.Errorf("unhandled charset: %s", charset)
		}
		if err := d.Decode(&task); err != nil {
			fmt.Println("unmarshal:", err)
			return CreateTaskResponse{}, err
		}
		return taskResponse, nil
	}
}

func (s *service) GetStatus(ctx context.Context, id string) {
	//TODO implement me
	panic("implement me")
}
