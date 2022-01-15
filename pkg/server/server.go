package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Selahattinn/picus-tcp-message/pkg/client"
	"github.com/Selahattinn/picus-tcp-message/pkg/command"
	"github.com/Selahattinn/picus-tcp-message/pkg/crypto"
)

// structure of server
type server struct {

	// map of clients connected to the server :
	// client name (key) & client (value)
	contacts map[string]*client.Client

	// channel on which server receives commands from clients
	commands chan command.Command
}

// function to instantiate new server
func newServer() *server {

	return &server{
		contacts: make(map[string]*client.Client),
		commands: make(chan command.Command),
	}
}

// function to run server
func (s *server) run() {

	log.Printf("running server...")

	// loop through incoming commands..
	for cmd := range s.commands {
		// based on the command id, execute desired functions
		switch cmd.ID {
		case command.CmdName:
			// update client name to input
			s.name(cmd.Client, cmd.Args[1])
		case command.CmdJoin:
			// update client contact to input
			s.join(cmd.Client, cmd.Args[1])
		case command.CmdList:
			// return list of users (clients) connected to the server
			s.list(cmd.Client)
		case command.CmdMsg:
			// send input to client contact
			s.msg(cmd.Client, cmd.Args)
		case command.CmdQuit:
			// quit chat system
			s.quit(cmd.Client)
		case command.CmdHelp:
			// return command list
			s.help(cmd.Client)
		}
	}
}

// function to instantiate new client :
// called when a new client joins the server
func (s *server) newClient(conn net.Conn) {

	// generate RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
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

	log.Printf("new client has joined : %s", conn.RemoteAddr().String())

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
		log.Printf("encrypting messages...")

		// send the message
		c.Msg(s.contacts[c.Contact], eMsg)
		log.Printf("sending message to %s", c.Contact)

	} else {

		// otherwise, prompt user to join to a user
		c.Msg(c, "no one hears you. follow below steps to get started :\n\n* use '/list' command to check, available users.\n* use '/join' command to select who you want to chat to.\n* use '/msg'  command to send message to selected user.\n")
	}

}

// function to exit from chat
func (s *server) quit(c *client.Client) {

	log.Printf("client has left the chat: %s", c.Conn.RemoteAddr().String())

	// remove user from server contact list
	_, ok := s.contacts[c.Name]
	if ok {
		delete(s.contacts, c.Name)
	}

	// pass message
	c.Msg(c, "skychat will miss you...")
	// close client connection
	c.Conn.Close()
}

// function to return command list
func (s *server) help(c *client.Client) {

	// pass message
	c.Msg(c, "Skychat : Chat Platform\n\n Usage : /<command> [arguments]\n\n* name : Specify your name.\n* list : List connected users.\n* join : Specify message recepient.\n* msg  : Send message to recepient.\n* quit : Exit Skychat.\n* help : List help commands.\n")

}
