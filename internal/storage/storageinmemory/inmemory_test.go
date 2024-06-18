package storageinmemory

import (
	"context"
	"testing"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageerrors"
	"github.com/stretchr/testify/assert"
)

func TestStorageInMemoryAdd(t *testing.T) {
	tests := map[string][]struct {
		messageId string
		taskId    string
	}{
		"simple": {{messageId: "1", taskId: "task-1"}},
		"a lot of add": {{messageId: "msg-1", taskId: "task-1"},
			{messageId: "msg-2", taskId: "task-2"},
			{messageId: "msg-3", taskId: "task-3"},
			{messageId: "msg-4", taskId: "task-4"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			for _, request := range tc {
				err := storage.AddElement(context.Background(), request.messageId, request.taskId)
				assert.Nil(t, err)
			}
			storage.Finish()
		})
	}
}

func TestGetMessageInMemory(t *testing.T) {
	tests := map[string]struct {
		taskId         string
		expectedResult string
		mockSetup      func(s *StorageMsgTasksInMemory)
	}{
		"simple": {expectedResult: "1", taskId: "task-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "1", "task-1")
			},
		},
		"a lot of elements": {expectedResult: "msg-1", taskId: "task-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-1", "task-1")
				s.AddElement(context.Background(), "msg-2", "task-2")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			tc.mockSetup(&storage)
			actualResult, err := storage.GetIdMessageByTask(context.Background(), tc.taskId)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actualResult)
			storage.Finish()
		})
	}
}

func TestGetTaskInMemory(t *testing.T) {
	tests := map[string]struct {
		msgId          string
		expectedResult string
		mockSetup      func(s *StorageMsgTasksInMemory)
	}{
		"simple": {expectedResult: "task-1", msgId: "msg-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-1", "task-1")
			},
		},
		"a lot of elements": {expectedResult: "task-1", msgId: "msg-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-1", "task-1")
				s.AddElement(context.Background(), "msg-2", "task-2")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			tc.mockSetup(&storage)
			actualResult, err := storage.GetIdTaskByMessage(context.Background(), tc.msgId)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actualResult)
			storage.Finish()
		})
	}
}

func TestGetTasksStorageInMemoryErrors(t *testing.T) {
	tests := map[string]struct {
		msgId       string
		expectedErr error
		mockSetup   func(s *StorageMsgTasksInMemory)
	}{
		"errorNotFound": {expectedErr: storageerrors.NewNotFoundError(), msgId: "msg-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-2", "task-1")
			},
		},
		"a lot of elements": {expectedErr: storageerrors.NewNotFoundError(), msgId: "msg-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-4", "task-1")
				s.AddElement(context.Background(), "msg-2", "task-2")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			tc.mockSetup(&storage)
			_, err := storage.GetIdTaskByMessage(context.Background(), tc.msgId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetMessagesStorageInMemoryErrors(t *testing.T) {
	tests := map[string]struct {
		taskId      string
		expectedErr error
		mockSetup   func(s *StorageMsgTasksInMemory)
	}{
		"errorNotFound": {expectedErr: storageerrors.NewNotFoundError(), taskId: "task-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-2", "task-2")
			},
		},
		"a lot of elements": {expectedErr: storageerrors.NewNotFoundError(), taskId: "t-1",
			mockSetup: func(s *StorageMsgTasksInMemory) {
				s.AddElement(context.Background(), "msg-4", "task-1")
				s.AddElement(context.Background(), "msg-2", "task-2")
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			storage := NewStorageMsgTasksInMemory()
			tc.mockSetup(&storage)
			_, err := storage.GetIdMessageByTask(context.Background(), tc.taskId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
