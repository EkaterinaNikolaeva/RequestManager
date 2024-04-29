package storagedb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type StorageMsgTasksDB struct {
	DB        *sql.DB
	tableName string
}

func NewStorageMsgTasksDB(login string, password string, host string, port string, name string, table string) (StorageMsgTasksDB, error) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		login,
		password,
		host,
		port,
		name)
	log.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return StorageMsgTasksDB{}, err
	}
	err = db.Ping()
	if err != nil {
		return StorageMsgTasksDB{}, err
	}
	return StorageMsgTasksDB{
		DB:        db,
		tableName: table,
	}, nil
}

func (s *StorageMsgTasksDB) AddElement(msgId string, taskId string, ctx context.Context) error {
	log.Printf("add message %s and task %s to db", msgId, taskId)
	query := fmt.Sprintf("INSERT into %s (idtask, idmessage) VALUES ('%s', '%s')", s.tableName, taskId, msgId)
	rows, err := s.DB.QueryContext(ctx, query)
	defer func() {
		err2 := rows.Close()
		if err2 != nil {
			return
		}
	}()
	return err
}

func (s *StorageMsgTasksDB) GetIdMessageByTask(taskId string, ctx context.Context) (string, bool, error) {
	log.Printf("get msg id by task id %s", taskId)
	query := fmt.Sprintf("select idmessage from %s where idtask='%s'", s.tableName, taskId)
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return "", false, err
	}
	defer rows.Close()
	for rows.Next() {
		var msgId string
		rows.Scan(&msgId)
		return msgId, true, nil
	}
	return "", false, nil
}

func (s *StorageMsgTasksDB) GetIdTaskByMessage(msgId string, ctx context.Context) (string, bool, error) {
	log.Printf("get task id by msg id %s", msgId)
	query := fmt.Sprintf("select idtask from %s where idmessage='%s'", s.tableName, msgId)
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return "", false, err
	}
	defer rows.Close()
	for rows.Next() {
		var taskId string
		rows.Scan(&taskId)
		taskId = strings.TrimSpace(taskId)
		return taskId, true, nil
	}
	return "", false, nil
}

func (s *StorageMsgTasksDB) Finish() {
	s.DB.Close()
}
