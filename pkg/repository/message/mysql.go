package message

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Selahattinn/picus-tcp-message/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type MySQLRepository struct {
	db *sql.DB
}

const (
	tableName = "messages"
)
const (
	initTableTemplate = `
	CREATE TABLE IF NOT EXISTS %s (
		id bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
		from_client TEXT NOT NULL,
		to_client TEXT NOT NULL,
		body TEXT NOT NULL,
		UNIQUE KEY id (id)
	  ) ENGINE=MyISAM  DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;	
`
)

func NewMySQLRepository(db *sql.DB) (*MySQLRepository, error) {
	tableInitCmd := fmt.Sprintf(initTableTemplate, tableName)
	_, err := db.Exec(tableInitCmd)

	if err != nil {
		return nil, fmt.Errorf("error init messages repository: %v", err)
	}

	return &MySQLRepository{
		db: db,
	}, nil
}

// GetAll returns all messages which is sended from a user
func (r *MySQLRepository) GetAll(from string) ([]model.Message, error) {
	q := "SELECT id, from_client, to_client, body FROM " + tableName + " where from_client=?"

	logrus.Debug("QUERY: ", q, from)
	res, err := r.db.Query(q, from)
	if err != nil {
		return nil, fmt.Errorf("error init message repository: %v", err)
	}
	var messages []model.Message
	for res.Next() {
		var message model.Message
		if err := res.Scan(&message.ID, &message.From, &message.To, &message.Text); err != nil {
			fmt.Println(err)
			return nil, err
		}
		messages = append(messages, message)

	}
	return messages, nil
}

// GetAll returns all messages which is sended from a user
func (r *MySQLRepository) GetAllToMe(from string) ([]model.Message, error) {
	q := "SELECT id, from_client, to_client, body FROM " + tableName + " where to_client=?"

	logrus.Debug("QUERY: ", q, from)
	res, err := r.db.Query(q, from)
	if err != nil {
		return nil, fmt.Errorf("error init message repository: %v", err)
	}
	var messages []model.Message
	for res.Next() {
		var message model.Message
		if err := res.Scan(&message.ID, &message.From, &message.To, &message.Text); err != nil {
			fmt.Println(err)
			return nil, err
		}
		messages = append(messages, message)

	}
	return messages, nil
}

// GetLast returns last X messages which is sended from a user
func (r *MySQLRepository) GetLast(from string, limit string) ([]model.Message, error) {
	q := "SELECT id, from_client, to_client, body FROM " + tableName + " where from_client=? ORDER BY id DESC LIMIT ?"

	logrus.Debug("QUERY: ", q, from)
	res, err := r.db.Query(q, from, limit)
	if err != nil {
		return nil, fmt.Errorf("error init message repository: %v", err)
	}
	var messages []model.Message
	for res.Next() {
		var message model.Message
		if err := res.Scan(&message.ID, &message.From, &message.To, &message.Text); err != nil {
			fmt.Println(err)
			return nil, err
		}
		messages = append(messages, message)

	}
	return messages, nil
}

// GetContains returns all messages which is contains a word
func (r *MySQLRepository) GetContains(from string, word string) ([]model.Message, error) {
	q := "SELECT id, from_client, to_client, body FROM " + tableName + " where from_client=?"

	logrus.Debug("QUERY: ", q)
	res, err := r.db.Query(q, from)
	if err != nil {
		return nil, fmt.Errorf("error init message repository: %v", err)
	}
	var messages []model.Message
	for res.Next() {
		var message model.Message
		if err := res.Scan(&message.ID, &message.From, &message.To, &message.Text); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if strings.Contains(message.Text, word) {
			messages = append(messages, message)
		}

	}
	return messages, nil
}

// Store returns an id which is ID of row
func (r *MySQLRepository) Store(message model.Message) (int64, error) {
	stmt, err := r.db.Prepare(`INSERT INTO ` + tableName + `(
		from_client,to_client,body)
		VALUES(
			?,?,?)`)
	if err != nil {
		return -1, err
	}

	defer stmt.Close()
	logrus.Debug("QUERY: ", stmt)
	res, err := stmt.Exec(
		message.From, message.To, message.Text)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
