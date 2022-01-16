package repository

import "github.com/Selahattinn/picus-tcp-message/pkg/repository/message"

// Repository defines the method for all operations related with repository
// Repository interface is composition of  Repository interfaces of imported packages.
type Repository interface {
	Shutdown()
	GetMessageRepository() message.Repository
}
