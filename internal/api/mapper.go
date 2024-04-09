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
	if task.StartMs > 0 {
		params = append(params, "")
		params = append(params, strconv.Itoa(task.StartMs))
	}
	if task.EndMs > 0 {
		params = append(params, strconv.Itoa(task.EndMs))
	}

	data := elecard.CreateTaskRequest{}
	data.Dispatcher = "1"
	data.SetValue.Name = "WatchFolderConfig"
	data.SetValue.Parameter.P = params

	return data
}
