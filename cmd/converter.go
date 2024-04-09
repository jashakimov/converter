package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jashakimov/converter/internal/api"
	"github.com/jashakimov/converter/internal/config"
	"github.com/jashakimov/converter/internal/service/elecard"
	"github.com/jashakimov/converter/internal/utils"
	"go.uber.org/zap"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.NewConfig(utils.GetConfigPath())

	gin.SetMode(gin.ReleaseMode)
	loggerProd, _ := zap.NewProduction()
	logger := loggerProd.Sugar()

	server := gin.New()
	server.Use(
		gin.Recovery(),
	)

	webSocketClient := NewWebSocketClient(cfg.ElecardWebSocket)
	elecardService := elecard.NewService(webSocketClient, logger)
	api.RegisterHandler(server, utils.NewValidator(), elecardService, logger)

	go func() {
		log.Println("Запущен сервер, порт", cfg.Port)
		if err := server.Run(":" + cfg.Port); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func NewWebSocketClient(host string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: host}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	return c
}
