package message

import "github.com/Selahattinn/picus-tcp-message/pkg/model"

type Reader interface {
	GetAll(from string) ([]model.Message, error)
	GetAllToMe(from string) ([]model.Message, error)
	GetLast(from string, limit string) ([]model.Message, error)
	GetContains(from string, word string) ([]model.Message, error)
}

type Writer interface {
	Store(message model.Message) (int64, error)
}

//Repository repository interface
type Repository interface {
	Reader
	Writer
}
