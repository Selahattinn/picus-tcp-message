package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"strings"

	"github.com/Selahattinn/picus-tcp-message/pkg/client"
	"github.com/Selahattinn/picus-tcp-message/pkg/crypto"
	"github.com/Selahattinn/picus-tcp-message/pkg/repository"
	"github.com/Selahattinn/picus-tcp-message/pkg/service"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ListenAddress string `yaml:"host"`

	Service *service.Config         `yaml:"service"`
	DB      *repository.MySQLConfig `yaml:"database"`
}

// structure of server
type server struct {

	// map of clients connected to the server :
	// client name (key) & client (value)
	contacts map[string]*client.Client

	// channel on which server receives commands from clients
	commands chan client.Command

	// Service Part
	Service service.Service

	// Yaml Config
	Config *Config
}

// function to instantiate new server
func NewServer(cfg *Config) *server {

	return &server{
		contacts: make(map[string]*client.Client),
		commands: make(chan client.Command),
		Config:   cfg,
	}
}

// function to run server
func (s *server) Run() {

	// Establish database connection
	repo, err := repository.NewMySQLRepository(s.Config.DB)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create mysql repository")
	}

	s.Service, err = service.NewProvider(s.Config.Service, repo)
	if err != nil {
		logrus.WithError(err).Fatal("Could not create service provider")
	}

	logrus.Info("running server...")

	// loop through incoming commands..
	for cmd := range s.commands {
		// based on the command id, execute desired functions
		switch cmd.ID {
		case client.CmdName:
			// update client name to input
			s.name(cmd.Client, cmd.Args[1])
		case client.CmdJoin:
			// update client contact to input
			s.join(cmd.Client, cmd.Args[1])
		case client.CmdList:
			// return list of users (clients) connected to the server
			s.list(cmd.Client)
		case client.CmdMsg:
			// send input to client contact
			s.msg(cmd.Client, cmd.Args)
		case client.CmdQuit:
			// quit chat system
			s.quit(cmd.Client)
		case client.CmdHelp:
			// return command list
			s.help(cmd.Client)
		}
	}
}

// function to instantiate new client :
// called when a new client joins the server
func (s *server) NewClient(conn net.Conn) {

	// generate RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.WithError(err).Fatal("RSA key generate fail in NewClient")
	}

	// instantiate client
	c := &client.Client{
		Conn:     conn,
		Name:     "anonymous",
		Commands: s.commands,
		Contact:  "",
		Private:  privateKey,
		Public:   privateKey.PublicKey,
	}
	logrus.Info("new client has joined : ", conn.RemoteAddr().String())

	// start reading for input ( this is a blocking call on a separte go routine )
	c.ReadInput()
}

// function to assign an identifer (name) to a newly created client
func (s *server) name(c *client.Client, name string) {

	// assign name to client
	c.Name = name

	// update server guest list i.e currently connected users (clients)
	s.contacts[name] = c

	// give user feedback message
	c.Msg(c, fmt.Sprintf("you will be known as %s", name))
}

// function to assign contact ( who a client is currently talkig to ) :
func (s *server) join(c *client.Client, contactName string) {

	// check if a user for given name exists on the server contacts map
	_, ok := s.contacts[contactName]

	// if so...
	if ok && contactName != "" {

		// update client contact ( this contact is who messages will be sent to )
		c.Contact = contactName
		// pass feedback
		c.Msg(c, fmt.Sprintf("You are now talking to :%s", c.Contact))

	} else {

		// otherwise, pass feedback
		c.Msg(c, "No such user exists. check available users again.")

	}
}

// function to display list of connected users :
// these clients are who you (a client) can join and then msg
func (s *server) list(c *client.Client) {

	var contacts []string

	// loop through available users
	for name := range s.contacts {

		// fetch all users except current client
		if name != c.Name {
			contacts = append(contacts, name)
		}

	}

	// pass message
	c.Msg(c, fmt.Sprintf("available users: %s", strings.Join(contacts, ", ")))
}

// function to pass a message to specified user (client)
func (s *server) msg(c *client.Client, args []string) {

	// check if a user for given name exists on the server contacts map
	_, ok := s.contacts[c.Contact]

	// is so...
	if ok && c.Contact != "" {

		// join the entire mesage
		msg := strings.Join(args[1:], " ")
		msg = c.Name + " : " + msg

		// fetch public key of recepient of message
		publicKey := s.contacts[c.Contact].Public

		// encrypt data
		eMsg := crypto.Encrypt(msg, publicKey)
		logrus.Info("encrypting messages... from client:", c.Name)

		// send the message
		c.Msg(s.contacts[c.Contact], eMsg)
		logrus.Info("sending message to ", c.Contact)

	} else {

		// otherwise, prompt user to join to a user
		c.Msg(c, "no one hears you. follow below steps to get started :\n\n* use '/list' command to check, available users.\n* use '/join' command to select who you want to chat to.\n* use '/msg'  command to send message to selected user.\n")
	}

}

// function to exit from chat
func (s *server) quit(c *client.Client) {
	logrus.Info("client has left the chat: ", c.Conn.RemoteAddr().String())

	// remove user from server contact list
	_, ok := s.contacts[c.Name]
	if ok {
		delete(s.contacts, c.Name)
	}

	// pass message
	c.Msg(c, "We will miss you...")
	// close client connection
	c.Conn.Close()
}

// function to return command list
func (s *server) help(c *client.Client) {

	// pass message
	c.Msg(c, "Picus Chat Platform\n\n Usage : /<command> [arguments]\n\n* name : Specify your name.\n* list : List connected users.\n* join : Specify message recepient.\n* msg  : Send message to recepient.\n* quit : Exit Chat App.\n* help : List help commands.\n")

}
