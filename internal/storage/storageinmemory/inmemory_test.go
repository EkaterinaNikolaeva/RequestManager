package storageinmemory

import (
	"context"
	"testing"

	errornotfound "github.com/EkaterinaNikolaeva/RequestManager/internal/storage/errors"
	"github.com/stretchr/testify/assert"
)

func TestStorageInMemory(t *testing.T) {
	type TypeRequest int
	const (
		add     TypeRequest = iota
		getMsg  TypeRequest = iota
		getTask TypeRequest = iota
	)
	tests := map[string][]struct {
		typeRequest    TypeRequest
		messageId      string
		taskId         string
		expectedResult string
		expectedErr    error
	}{
		"simple": {{typeRequest: add, messageId: "1", taskId: "task-1", expectedErr: nil},
			{typeRequest: getMsg, taskId: "task-1", expectedResult: "1", expectedErr: nil},
			{typeRequest: getTask, messageId: "1", expectedResult: "task-1", expectedErr: nil},
		},
		"no elements": {{typeRequest: getMsg, taskId: "task-1", expectedErr: errornotfound.NewNotFoundError()},
			{typeRequest: getTask, messageId: "1", expectedErr: errornotfound.NewNotFoundError()},
			{typeRequest: getTask, messageId: "other-id-msg", expectedErr: errornotfound.NewNotFoundError()},
			{typeRequest: getMsg, taskId: "othertaskid", expectedErr: errornotfound.NewNotFoundError()},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			for _, request := range tc {
				if request.typeRequest == add {
					err := storage.AddElement(context.Background(), request.messageId, request.taskId)
					assert.Equal(t, request.expectedErr, err)
				} else if request.typeRequest == getMsg {
					res, err := storage.GetIdMessageByTask(context.Background(), request.taskId)
					assert.Equal(t, request.expectedResult, res)
					assert.Equal(t, request.expectedErr, err)
				} else {
					res, err := storage.GetIdTaskByMessage(context.Background(), request.messageId)
					assert.Equal(t, request.expectedResult, res)
					assert.Equal(t, request.expectedErr, err)
				}
			}
		})
	}
}
