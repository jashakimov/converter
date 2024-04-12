package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jashakimov/converter/internal/api"
	"github.com/jashakimov/converter/internal/config"
	"github.com/jashakimov/converter/internal/service/elecard"
	"github.com/jashakimov/converter/internal/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc

	cfg := config.NewConfig(utils.GetConfigPath())
	logger := NewLogger(cfg)
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(
		gin.Recovery(),
	)

	logger.Info("Подключение к ", cfg.ElecardWebSocket)
	webSocketClient := NewWebSocketClient(cfg.ElecardWebSocket)
	elecardService := elecard.NewService(webSocketClient, logger)
	api.RegisterHandler(
		server,
		utils.NewValidator(),
		elecardService,
		logger,
		time.Duration(cfg.TimeoutSec)*time.Second,
	)

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

func NewLogger(cfg config.Config) *zap.SugaredLogger {
	rotator, err := rotatelogs.New(
		cfg.LogPath+"/%Y/%m/%d.log",
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		panic(err)
	}
	encoderConfig := map[string]string{
		"levelEncoder": "capital",
		"timeKey":      "date",
		"timeEncoder":  "iso8601",
		"messageKey":   "message",
		"levelKey":     "level",
	}
	data, _ := json.Marshal(encoderConfig)
	var encCfg zapcore.EncoderConfig
	if err := json.Unmarshal(data, &encCfg); err != nil {
		panic(err)
	}
	w := zapcore.AddSync(rotator)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		w,
		zap.InfoLevel)
	return zap.New(core).Sugar()
}
