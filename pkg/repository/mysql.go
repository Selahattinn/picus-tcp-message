package repository

import (
	"database/sql"

	"github.com/Selahattinn/picus-tcp-message/pkg/repository/message"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL Repository defines the MySQL implementation of Repository interface
type MySQLRepository struct {
	cfg               *MySQLConfig
	db                *sql.DB
	messageRepository message.Repository
}

// MySQLConfig defines the MySQL Repository configuration
type MySQLConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// dbConn opens connection with MySQL driver
func dbConn(cfg *MySQLConfig) (*sql.DB, error) {

	dbDriver := "mysql"    // Database driver
	dbUser := cfg.Username // Mysql username
	dbPass := cfg.Password // Mysql password
	dbName := cfg.DBName   // Mysql schema
	addr := cfg.Addr

	// Realize the connection with mysql driver
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+addr+")/")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, err
	}

	db.Close()

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	// Return db object to be used by other functions
	return db, nil
}

// NewMySQLRepository creates a new MySQL Repository
func NewMySQLRepository(cfg *MySQLConfig) (*MySQLRepository, error) {
	db, err := dbConn(cfg)
	if err != nil {
		return nil, err
	}

	messageRepository, err := message.NewMySQLRepository(db)
	if err != nil {
		return nil, err
	}
	return &MySQLRepository{
		cfg:               cfg,
		db:                db,
		messageRepository: messageRepository,
	}, nil
}

// GetMessageRepository returns the message repository
func (r *MySQLRepository) GetMessageRepository() message.Repository {
	return r.messageRepository
}

// Shutdown closes the database connection
func (r *MySQLRepository) Shutdown() {
	r.db.Close()
}
