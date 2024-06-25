package storagepostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/EkaterinaNikolaeva/RequestManager/internal/storage/storageerrors"
	_ "github.com/lib/pq"
)

type StorageMsgTasksDB struct {
	DB        *sql.DB
	tableName string
}

func NewStorageMsgTasksDB(ctx context.Context, login string, password string, host string, port string, name string, table string, sqlOpen func(string, string) (*sql.DB, error)) (StorageMsgTasksDB, error) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		login,
		password,
		host,
		port,
		name)
	db, err := sqlOpen("postgres", connStr)
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
	query := "INSERT INTO " + s.tableName + " (idtask, idmessage) VALUES ($1, $2)"
	_, err := s.DB.ExecContext(ctx, query, taskId, msgId)
	return err
}

func (s *StorageMsgTasksDB) GetIdMessageByTask(ctx context.Context, taskId string) (string, error) {
	log.Printf("get msg id by task id %s", taskId)
	query := fmt.Sprintf("select idmessage from %s where idtask=$1", s.tableName)
	rows, err := s.DB.QueryContext(ctx, query, taskId)
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
	return "", storageerrors.NewNotFoundError()
}

func (s *StorageMsgTasksDB) GetIdTaskByMessage(ctx context.Context, msgId string) (string, error) {
	query := fmt.Sprintf("select idtask from %s where idmessage=$1", s.tableName)
	rows, err := s.DB.QueryContext(ctx, query, msgId)
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
	return "", storageerrors.NewNotFoundError()
}

func (s *StorageMsgTasksDB) Finish() {
	s.DB.Close()
}
