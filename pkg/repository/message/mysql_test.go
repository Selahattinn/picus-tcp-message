package message

import (
	"database/sql"
	"log"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Selahattinn/picus-tcp-message/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var m = &model.Message{
	ID:   int64(1),
	From: "Test",
	To:   "Test2",
	Text: "Test Text",
}

var wrongM = &model.Message{
	ID:   int64(1),
	From: "Test3",
	To:   "Test4",
	Text: "Test Text2",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalln(err)
	}
	repo := &MySQLRepository{db: db}
	query := "SELECT id, from_client, to_client, body FROM messages where from_client=?"

	rows := sqlmock.NewRows([]string{"id", "from_client", "to_client", "body"}).
		AddRow(m.ID, m.From, m.To, m.Text)

	mock.ExpectQuery(query).WithArgs(m.From).WillReturnRows(rows)

	messages, err := repo.GetAll(m.From)
	assert.NotNil(t, messages)
	assert.NoError(t, err)

}

func TestMySQLRepository_GetAllToMe(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalln(err)
	}
	repo := &MySQLRepository{db: db}
	query := "SELECT id, from_client, to_client, body FROM messages where to_client=?"

	rows := sqlmock.NewRows([]string{"id", "from_client", "to_client", "body"}).
		AddRow(m.ID, m.From, m.To, m.Text)

	mock.ExpectQuery(query).WithArgs(m.From).WillReturnRows(rows)

	messages, err := repo.GetAllToMe(m.From)
	assert.NotNil(t, messages)
	assert.NoError(t, err)
}

func TestMySQLRepository_GetLast(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalln(err)
	}
	repo := &MySQLRepository{db: db}
	query := "SELECT id, from_client, to_client, body FROM messages where from_client=? ORDER BY id DESC LIMIT ?"

	rows := sqlmock.NewRows([]string{"id", "from_client", "to_client", "body"}).
		AddRow(m.ID, m.From, m.To, m.Text)

	mock.ExpectQuery(query).WithArgs(m.From, "2").WillReturnRows(rows)

	message, err := repo.GetLast(m.From, "2")
	assert.NotNil(t, message)
	assert.NoError(t, err)
}

func TestMySQLRepository_GetContains(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalln(err)
	}
	repo := &MySQLRepository{db: db}
	query := "SELECT id, from_client, to_client, body FROM messages where from_client=?"

	rows := sqlmock.NewRows([]string{"id", "from_client", "to_client", "body"}).
		AddRow(m.ID, m.From, m.To, m.Text)

	mock.ExpectQuery(query).WithArgs(m.From).WillReturnRows(rows)

	messages, err := repo.GetContains(m.From, "Test")
	if err != nil {
		log.Fatalln(err)
	}
	for _, message := range messages {
		if strings.Contains(message.Text, "Test") {
			return
		}
	}
	log.Fatal("Get COntains not working")

	assert.NotNil(t, messages)
	assert.NoError(t, err)
}
