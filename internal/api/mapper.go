package api

import (
	"github.com/jashakimov/converter/internal/service/elecard"
	"strconv"
)

type mapper struct{}

func (m *mapper) NewCreateTaskRequest(task Task) elecard.CreateTaskRequest {
	params := []string{
		"Put",
		task.ID,
		task.Path,
	}
	if task.StartMs > 0 && task.EndMs > 0 {
		params = append(params, "")
		params = append(params, strconv.Itoa(task.StartMs))
		params = append(params, strconv.Itoa(task.EndMs))
	}

	taskRequest := elecard.CreateTaskRequest{}
	taskRequest.Dispatcher = "1"
	taskRequest.SetValue.Name = "WatchFolderConfig"
	taskRequest.SetValue.Parameter.P = params

	return taskRequest
}

func (m *mapper) NewGetStatusRequest(id string) elecard.GetStatusRequest {
	statusRequest := elecard.GetStatusRequest{}
	statusRequest.Dispatcher = "1"
	statusRequest.GetValue.Name = "WatchFolderConfig"
	statusRequest.GetValue.Parameter.P = []string{
		"Queue",
		id,
	}

	return statusRequest
}
