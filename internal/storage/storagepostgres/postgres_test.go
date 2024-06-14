package storagepostgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageerrors"
	"github.com/tj/assert"
)

func TestNewStorageMsgTasksDB(t *testing.T) {
	tests := map[string]struct {
		login     string
		password  string
		host      string
		port      string
		name      string
		table     string
		mockSetup func(sqlmock.Sqlmock)
	}{
		"success": {
			login:    "user",
			password: "pass",
			host:     "localhost",
			port:     "5432",
			name:     "test",
			table:    "table",
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
			sqlOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
				return db, nil
			}
			storage, err := NewStorageMsgTasksDB(ctx, tc.login, tc.password, tc.host, tc.port, tc.name, tc.table, sqlOpen)
			assert.NoError(t, err)
			assert.Equal(t, tc.table, storage.tableName)
			assert.NotNil(t, storage.DB)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddElement(t *testing.T) {
	tests := map[string]struct {
		table     string
		msgId     string
		taskId    string
		mockSetup func(sqlmock.Sqlmock)
	}{
		"simple": {
			table:  "test",
			msgId:  "msg-1",
			taskId: "task-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO test (idtask, idmessage) VALUES ($1, $2)").
					WithArgs("task-1", "msg-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer db.Close()
			tc.mockSetup(mock)
			storage := StorageMsgTasksDB{
				DB:        db,
				tableName: tc.table,
			}
			err = storage.AddElement(context.Background(), tc.msgId, tc.taskId)
			if err != nil {
				t.Errorf("TestAddElement() error = %v", err)
				return
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddElementErrors(t *testing.T) {
	tests := map[string]struct {
		table     string
		msgId     string
		taskId    string
		mockSetup func(sqlmock.Sqlmock)
	}{
		"return error": {
			table:  "test",
			msgId:  "msg-1",
			taskId: "task-1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO test (idtask, idmessage) VALUES ($1, $2)").
					WithArgs("task-1", "msg-1").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			assert.NoError(t, err)
			defer db.Close()
			tc.mockSetup(mock)
			storage := StorageMsgTasksDB{
				DB:        db,
				tableName: tc.table,
			}
			err = storage.AddElement(context.Background(), tc.msgId, tc.taskId)
			assert.Error(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestGetIdTaskByMessage(t *testing.T) {
	tests := map[string]struct {
		table        string
		msgId        string
		expectedTask string
		mockSetup    func(sqlmock.Sqlmock)
	}{
		"correct": {
			table:        "test",
			msgId:        "msg",
			expectedTask: "task",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"idtask"}).AddRow("task")
				mock.ExpectQuery("select idtask from test where idmessage=\\$1").
					WithArgs("msg").
					WillReturnRows(rows)
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
			if err != nil {
				t.Errorf("TestAddElement() error = %v", err)
				return
			}
			assert.Equal(t, tc.expectedTask, msg)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetIdTaskByMessageErrors(t *testing.T) {
	tests := map[string]struct {
		table       string
		msgId       string
		expectedErr error
		mockSetup   func(sqlmock.Sqlmock)
	}{
		"no such task": {
			table:       "test",
			msgId:       "msg",
			expectedErr: storageerrors.NewNotFoundError(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select idtask from test where idmessage=\\$1").
					WithArgs("msg").
					WillReturnError(storageerrors.NewNotFoundError())
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
			_, err = storage.GetIdTaskByMessage(context.Background(), tc.msgId)
			if err == nil {
				t.Errorf("TestGetIdTaskByMessageErrors() error = %v", err)
				return
			}
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetIdMessageByTask(t *testing.T) {
	tests := map[string]struct {
		table       string
		taskId      string
		expectedMsg string
		mockSetup   func(sqlmock.Sqlmock)
	}{
		"correct": {
			table:       "test",
			taskId:      "task",
			expectedMsg: "msg",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"idmessage"}).AddRow("msg")
				mock.ExpectQuery("select idmessage from test where idtask=\\$1").
					WithArgs("task").
					WillReturnRows(rows)
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
			if err != nil {
				t.Errorf("TestGetIdMessageByTask() error = %v", err)
				return
			}
			assert.Equal(t, tc.expectedMsg, msg)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetIdMessageByTaskErrors(t *testing.T) {
	tests := map[string]struct {
		table       string
		taskId      string
		expectedErr error
		mockSetup   func(sqlmock.Sqlmock)
	}{
		"no such msg": {
			table:       "test",
			taskId:      "task",
			expectedErr: storageerrors.NewNotFoundError(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select idmessage from test where idtask=\\$1").
					WithArgs("task").
					WillReturnError(storageerrors.NewNotFoundError())
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
			_, err = storage.GetIdMessageByTask(context.Background(), tc.taskId)
			if err == nil {
				t.Errorf("TestGetIdMessageByTaskErrors() not error")
				return
			}
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
