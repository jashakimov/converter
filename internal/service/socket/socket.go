package socket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/url"
)

type Service interface {
	Send(request []byte) ([]byte, error)
}

type service struct {
	lg         *zap.SugaredLogger
	connString string
}

func NewService(connString string, lg *zap.SugaredLogger) Service {
	s := service{
		connString: connString,
		lg:         lg,
	}
	return &s
}

func (s *service) Send(request []byte) ([]byte, error) {
	u := url.URL{Scheme: "ws", Host: s.connString}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		s.lg.Error("Ошибка подключения к веб-сокету:", err.Error())
		return nil, err
	}
	s.lg.Info("Подключение к веб-сокету:", s.connString)
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
	}()

	if err := conn.WriteMessage(websocket.TextMessage, request); err != nil {
		s.lg.Error("Ошибка записи в веб-сокет:", err.Error())
		return nil, err
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.lg.Error("Ошибка чтения из веб-сокета:", err.Error())
			return nil, err
		}
		return msg, nil
	}
}
