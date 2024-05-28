package storagepostgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	errornotfound "github.com/EkaterinaNikolaeva/RequestManager/internal/storage/errors"
	"github.com/tj/assert"
)

func TestNewStorageMsgTasksDB(t *testing.T) {
	tests := map[string]struct {
		login         string
		password      string
		host          string
		port          string
		name          string
		table         string
		expectedError bool
		mockSetup     func(sqlmock.Sqlmock)
	}{
		"success": {
			login:         "user",
			password:      "pass",
			host:          "localhost",
			port:          "5432",
			name:          "test",
			table:         "table",
			expectedError: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing()
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			tc.mockSetup(mock)
			ctx := context.Background()
			oldOpen := sqlOpen
			sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
				return db, nil
			}
			defer func() { sqlOpen = oldOpen }()
			storage, err := NewStorageMsgTasksDB(ctx, tc.login, tc.password, tc.host, tc.port, tc.name, tc.table)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.table, storage.tableName)
				assert.NotNil(t, storage.DB)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddElement(t *testing.T) {
	tests := map[string]struct {
		table     string
		msgId     string
		taskId    string
		wantErr   bool
		mockSetup func(sqlmock.Sqlmock)
	}{
		"simple": {
			table:   "test",
			msgId:   "msg-1",
			taskId:  "task-1",
			wantErr: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO test \\(idtask, idmessage\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("task-1", "msg-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"return error": {
			table:   "test",
			msgId:   "msg-1",
			taskId:  "task-1",
			wantErr: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO test \\(idtask, idmessage\\) VALUES \\(\\$1, \\$2\\)").
					WithArgs("task-1", "msg-1").
					WillReturnError(nil)
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			tc.mockSetup(mock)
			storage := StorageMsgTasksDB{
				DB:        db,
				tableName: tc.table,
			}
			err = storage.AddElement(context.Background(), tc.msgId, tc.taskId)
			if (err != nil) != tc.wantErr {
				t.Errorf("TestAddElement() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetIdTaskByMessage(t *testing.T) {
	tests := map[string]struct {
		table        string
		msgId        string
		wantErr      bool
		expectedErr  error
		expectedTask string
		mockSetup    func(sqlmock.Sqlmock)
	}{
		"correct": {
			table:        "test",
			msgId:        "msg",
			wantErr:      false,
			expectedTask: "task",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"idtask"}).AddRow("task")
				mock.ExpectQuery("select idtask from test where idmessage=\\$1").
					WithArgs("msg").
					WillReturnRows(rows)
			},
		},
		"no such task": {
			table:       "test",
			msgId:       "msg",
			wantErr:     true,
			expectedErr: errornotfound.NewNotFoundError(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select idtask from test where idmessage=\\$1").
					WithArgs("msg").
					WillReturnError(errornotfound.NewNotFoundError())
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)
			storage := StorageMsgTasksDB{
				DB:        db,
				tableName: tc.table,
			}
			msg, err := storage.GetIdTaskByMessage(context.Background(), tc.msgId)
			if (err != nil) != tc.wantErr {
				t.Errorf("TestAddElement() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.expectedTask, msg)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetIdMessageByTask(t *testing.T) {
	tests := map[string]struct {
		table       string
		taskId      string
		wantErr     bool
		expectedErr error
		expectedMsg string
		mockSetup   func(sqlmock.Sqlmock)
	}{
		"correct": {
			table:       "test",
			taskId:      "task",
			wantErr:     false,
			expectedMsg: "msg",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"idmessage"}).AddRow("msg")
				mock.ExpectQuery("select idmessage from test where idtask=\\$1").
					WithArgs("task").
					WillReturnRows(rows)
			},
		},
		"no such task": {
			table:       "test",
			taskId:      "task",
			wantErr:     true,
			expectedErr: errornotfound.NewNotFoundError(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select idmessage from test where idtask=\\$1").
					WithArgs("task").
					WillReturnError(errornotfound.NewNotFoundError())
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)
			storage := StorageMsgTasksDB{
				DB:        db,
				tableName: tc.table,
			}
			msg, err := storage.GetIdMessageByTask(context.Background(), tc.taskId)
			if (err != nil) != tc.wantErr {
				t.Errorf("TestAddElement() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.expectedMsg, msg)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
