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
	conn       *websocket.Conn
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
	s.connect()
	defer s.conn.Close()

	//c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	if err := s.conn.WriteMessage(websocket.TextMessage, request); err != nil {
		s.lg.Error("Ошибка записи в веб-сокет:", err.Error())
		return nil, err
	}
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			s.lg.Error("Ошибка чтения из веб-сокета:", err.Error())
			return nil, err
		}
		return msg, nil
	}
}

func (s *service) connect() {
	s.lg.Info("Поключаемся к сокету")
	u := url.URL{Scheme: "ws", Host: s.connString}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	s.conn = c
}
