package message

import (
	"github.com/Selahattinn/picus-tcp-message/pkg/model"
	"github.com/Selahattinn/picus-tcp-message/pkg/repository"
)

type Service struct {
	repository repository.Repository
}

func NewService(repo repository.Repository) (*Service, error) {
	return &Service{
		repository: repo,
	}, nil
}

// GetAllMessages returns all messages which is sended
func (s *Service) GetAllMessages(from_client string) ([]model.Message, error) {
	messages, err := s.repository.GetMessageRepository().GetAll(from_client)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetAllMessagesToMe returns all messages which is recived
func (s *Service) GetAllMessagesToMe(from_client string) ([]model.Message, error) {
	messages, err := s.repository.GetMessageRepository().GetAllToMe(from_client)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetLast returns last X messages
func (s *Service) GetLast(from_client string, limit string) ([]model.Message, error) {
	messages, err := s.repository.GetMessageRepository().GetLast(from_client, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetLast returns all messages which is containts a word
func (s *Service) GetContains(from_client string, word string) ([]model.Message, error) {
	messages, err := s.repository.GetMessageRepository().GetContains(from_client, word)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// StoreMessage for storing a message
func (s *Service) StoreMessage(message model.Message) (int64, error) {
	id, err := s.repository.GetMessageRepository().Store(message)
	if err != nil {
		return -1, err
	}
	return id, nil

}
