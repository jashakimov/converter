package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jashakimov/converter/internal/service/elecard"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type api struct {
	lg             *zap.SugaredLogger
	elecardService elecard.Service
	val            *validator.Validate
	m              mapper
}

func RegisterHandler(
	router *gin.Engine,
	val *validator.Validate,
	elecardService elecard.Service,
	lg *zap.SugaredLogger,
) {
	api := api{
		lg:             lg,
		elecardService: elecardService,
		val:            val,
		m:              mapper{},
	}

	router.POST("/tasks", api.CreateTask)
}

func (a *api) CreateTask(ctx *gin.Context) {
	ctxWithT, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var request Task
	if err := ctx.BindJSON(&request); err != nil {
		a.lg.Error(err)
		return
	}

	if err := a.val.Struct(request); err != nil {
		a.lg.Error(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	taskResponse, err := a.elecardService.CreateTask(ctxWithT, a.m.NewCreateTaskRequest(request))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, taskResponse)
}
