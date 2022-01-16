package service

import (
	"github.com/Selahattinn/picus-tcp-message/pkg/repository"
	"github.com/Selahattinn/picus-tcp-message/pkg/service/message"
)

type Provider struct {
	cfg            *Config
	repository     repository.Repository
	messageService *message.Service
}

func NewProvider(cfg *Config, repo repository.Repository) (*Provider, error) {
	messageService, err := message.NewService(repo)
	if err != nil {
		return nil, err
	}
	return &Provider{
		cfg:            cfg,
		repository:     repo,
		messageService: messageService,
	}, nil
}

func (p *Provider) GetConfig() *Config {
	return p.cfg
}
func (p *Provider) GetMessageService() *message.Service {
	return p.messageService
}
func (p *Provider) Shutdown() {
	p.repository.Shutdown()
}
