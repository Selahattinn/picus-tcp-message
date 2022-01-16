package service

import "github.com/Selahattinn/picus-tcp-message/pkg/service/message"

type Config struct{}

type Service interface {
	GetConfig() *Config
	GetMessageService() *message.Service
	Shutdown()
}
