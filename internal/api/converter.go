package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jashakimov/converter/internal/config"
	"github.com/jashakimov/converter/internal/service/elecard"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type api struct {
	lg             *zap.SugaredLogger
	elecardService elecard.Service
	val            *validator.Validate
	m              mapper
	timeout        time.Duration
	requestDelay   time.Duration
}

func RegisterHandler(
	router *gin.Engine,
	val *validator.Validate,
	elecardService elecard.Service,
	lg *zap.SugaredLogger,
	cfg config.Config,
) {
	api := api{
		lg:             lg,
		elecardService: elecardService,
		val:            val,
		m:              mapper{},
		timeout:        time.Duration(cfg.TimeoutSec) * time.Second,
		requestDelay:   time.Duration(cfg.RequestDelaySec) * time.Second,
	}

	router.POST("/tasks", api.CreateTask)
}

func (a *api) CreateTask(ctx *gin.Context) {
	ctxWithTime, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	chanResult := make(chan ChanResult)
	go func() {
		startMs, _ := strconv.Atoi(ctx.Query("startMs"))
		endMs, _ := strconv.Atoi(ctx.Query("endMs"))

		request := Task{
			Path:    ctx.Query("path"),
			ID:      ctx.Query("id"),
			StartMs: startMs,
			EndMs:   endMs,
		}

		//if err := ctx.BindJSON(&request); err != nil {
		//	a.lg.Error(err)
		//	chanResult <- ChanResult{
		//		Error: err,
		//	}
		//	return
		//}
		a.lg.Info("Пришел запрос на вход:", ctx.Request.URL.RawQuery)

		//if err := a.val.Struct(request); err != nil {
		//	a.lg.Error(err)
		//	chanResult <- ChanResult{
		//		Error: err,
		//	}
		//	return
		//}

		taskResponse, err := a.elecardService.CreateTask(ctxWithTime, a.m.NewCreateTaskRequest(request))
		if err != nil {
			chanResult <- ChanResult{
				Error: err,
			}
			return
		}

		time.Sleep(a.requestDelay)

		splitted := strings.Split(request.Path, "\\")
		fileName := splitted[len(splitted)-1]
		status, err := a.elecardService.GetStatus(
			ctxWithTime,
			a.m.NewGetStatusRequest(taskResponse.SetValue.RetVal.WatchFolder.ID),
			fileName,
			a.requestDelay,
		)
		if err != nil {
			chanResult <- ChanResult{
				Error: err,
			}
			return
		}
		chanResult <- ChanResult{
			Error:   nil,
			Message: status,
		}
	}()

	select {
	case <-ctxWithTime.Done():
		ctx.String(http.StatusGatewayTimeout, "%s", "Timeout is expired")
		return
	case status := <-chanResult:
		if status.Error == nil {
			ctx.String(http.StatusOK, "%s", status.Message)
		} else {
			ctx.String(http.StatusBadRequest, "%s", status.Error.Error())
		}
		return
	}
}
