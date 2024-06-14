package storageinmemory

import (
	"context"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageerrors"
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
		"no elements": {{typeRequest: getMsg, taskId: "task-1", expectedErr: storageerrors.NewNotFoundError()},
			{typeRequest: getTask, messageId: "1", expectedErr: storageerrors.NewNotFoundError()},
			{typeRequest: getTask, messageId: "other-id-msg", expectedErr: storageerrors.NewNotFoundError()},
			{typeRequest: getMsg, taskId: "othertaskid", expectedErr: storageerrors.NewNotFoundError()},
		},
		"a lot of add": {{typeRequest: add, messageId: "msg-1", taskId: "task-1", expectedErr: nil},
			{typeRequest: add, messageId: "msg-2", taskId: "task-2", expectedErr: nil},
			{typeRequest: add, messageId: "msg-3", taskId: "task-3", expectedErr: nil},
			{typeRequest: getMsg, taskId: "task-1", expectedResult: "msg-1"},
			{typeRequest: getTask, messageId: "msg-2", expectedResult: "task-2"},
			{typeRequest: add, messageId: "msg-4", taskId: "task-4", expectedErr: nil},
			{typeRequest: getTask, messageId: "msg-4", expectedResult: "task-4"},
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
