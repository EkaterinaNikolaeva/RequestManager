package storagedb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	errornotfound "github.com/EkaterinaNikolaeva/RequestManager/internal/storage/errors"
	_ "github.com/lib/pq"
)

type StorageMsgTasksDB struct {
	DB        *sql.DB
	tableName string
}

func NewStorageMsgTasksDB(ctx context.Context, login string, password string, host string, port string, name string, table string) (StorageMsgTasksDB, error) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		login,
		password,
		host,
		port,
		name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return StorageMsgTasksDB{}, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return StorageMsgTasksDB{}, err
	}
	return StorageMsgTasksDB{
		DB:        db,
		tableName: table,
	}, nil
}

func (s *StorageMsgTasksDB) AddElement(ctx context.Context, msgId string, taskId string) error {
	log.Printf("add message %s and task %s to db", msgId, taskId)
	query := fmt.Sprintf("INSERT into %s (idtask, idmessage) VALUES ('%s', '%s')", s.tableName, taskId, msgId)
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	err = rows.Close()
	return err
}

func (s *StorageMsgTasksDB) GetIdMessageByTask(ctx context.Context, taskId string) (string, error) {
	log.Printf("get msg id by task id %s", taskId)
	query := fmt.Sprintf("select idmessage from %s where idtask='%s'", s.tableName, taskId)
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var msgId string
		err = rows.Scan(&msgId)
		if err != nil {
			return "", err
		}
		return msgId, nil
	}
	return "", errornotfound.NewNotFoundError()
}

func (s *StorageMsgTasksDB) GetIdTaskByMessage(ctx context.Context, msgId string) (string, error) {
	query := fmt.Sprintf("select idtask from %s where idmessage='%s'", s.tableName, msgId)
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var taskId string
		rows.Scan(&taskId)
		taskId = strings.TrimSpace(taskId)
		return taskId, nil
	}
	return "", errornotfound.NewNotFoundError()
}

func (s *StorageMsgTasksDB) Finish() {
	s.DB.Close()
}
