package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jashakimov/converter/internal/service/elecard"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"time"
)

type api struct {
	lg             *zap.SugaredLogger
	elecardService elecard.Service
	val            *validator.Validate
	m              mapper
	timeout        time.Duration
}

func RegisterHandler(
	router *gin.Engine,
	val *validator.Validate,
	elecardService elecard.Service,
	lg *zap.SugaredLogger,
	timeout time.Duration,
) {
	api := api{
		lg:             lg,
		elecardService: elecardService,
		val:            val,
		m:              mapper{},
		timeout:        timeout,
	}

	router.POST("/tasks", api.CreateTask)
}

func (a *api) CreateTask(ctx *gin.Context) {
	ctxWithTime, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	chanResult := make(chan ChanResult)
	go func() {
		var request Task
		if err := ctx.BindJSON(&request); err != nil {
			a.lg.Error(err)
			chanResult <- ChanResult{
				Error: err,
			}
			return
		}
		a.lg.Info("Пришел на вход:", request)

		if err := a.val.Struct(request); err != nil {
			a.lg.Error(err)
			chanResult <- ChanResult{
				Error: err,
			}
			return
		}

		taskResponse, err := a.elecardService.CreateTask(ctxWithTime, a.m.NewCreateTaskRequest(request))
		if err != nil {
			chanResult <- ChanResult{
				Error: err,
			}
			return
		}

		status, err := a.elecardService.GetStatus(
			ctxWithTime,
			a.m.NewGetStatusRequest(taskResponse.SetValue.RetVal.WatchFolder.ID),
			filepath.Base(request.Path),
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
