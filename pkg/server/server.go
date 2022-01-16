package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Selahattinn/picus-tcp-message/pkg/client"
	"github.com/Selahattinn/picus-tcp-message/pkg/crypto"
	"github.com/Selahattinn/picus-tcp-message/pkg/model"
	"github.com/Selahattinn/picus-tcp-message/pkg/repository"
	"github.com/Selahattinn/picus-tcp-message/pkg/service"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// Host adress which server run
	ListenAddress string `yaml:"host"`

	// Service configs
	Service *service.Config `yaml:"service"`
	// DB configs
	DB *repository.MySQLConfig `yaml:"database"`
}
type CmdStruct struct {
	Type  string
	Value string
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
			s.name(cmd.Client, cmd.Args)
		case client.CmdJoin:
			// update client contact to input
			s.join(cmd.Client, cmd.Args)
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
		case client.CmdGetMessageFromMe:
			// return Messages with limited from itself
			s.getMessageFromMe(cmd.Client, cmd.Args)
		case client.CmdGetMessageToMe:
			// return Messages with limited from itself
			s.getMessageToMe(cmd.Client, cmd.Args)
		case client.CmdGetLast:
			// return Messages with limited from itself
			s.getLastMassge(cmd.Client, cmd.Args)
		case client.CmdGetContains:
			// return Messages with limited from itself
			s.getContains(cmd.Client, cmd.Args)
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
func (s *server) name(c *client.Client, args []string) {
	if len(args) > 2 || len(args) < 2 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/name Selahattin")
		return
	}
	// assign name to client
	c.Name = args[1]

	// update server guest list i.e currently connected users (clients)
	// Control for client name
	// Client name can not be equal to any clients name
	for _, y := range s.contacts {
		if y.Name == c.Name {
			c.Msg(c, "There is a user which is used for this name. Please choose another name")
			return
		}
	}
	s.contacts[args[1]] = c

	// give user feedback message
	c.Msg(c, fmt.Sprintf("you will be known as %s", args[1]))
}

// function to assign contact ( who a client is currently talkig to ) :
func (s *server) join(c *client.Client, args []string) {
	if len(args) > 2 || len(args) < 2 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/join Selahattin")
		return
	}
	// check if a user for given name exists on the server contacts map
	_, ok := s.contacts[args[1]]

	// if so...
	if ok && args[1] != "" {

		// update client contact ( this contact is who messages will be sent to )
		c.Contact = args[1]
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
		message := &model.Message{
			From: c.Name,
			To:   c.Contact,
			Text: strings.Join(args[1:], " "),
		}

		_, err := s.Service.GetMessageService().StoreMessage(*message)
		if err != nil {
			logrus.WithError(err).Info("Message not saved to db")
		}
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
	c.Msg(c, "Picus Chat Platform\n\n Usage : /<command> [arguments]\n\n* name : Specify your name.\n* list : List connected users.\n* join : Specify message recepient.\n* msg  : Send message to recepient.\n* quit : Exit Chat App.\n* help : List help commands.\n* get-last : List of last sended messages.\n* get-contains : List of messages which is include this word.\n* get-m-to-me : lists all messages sent to me.\n* get-m-from-me : Lists all the messages I've sent.\n")

}

// For write to msg last X messages whic is sendend from me
func (s *server) getMessageFromMe(c *client.Client, args []string) {
	if len(args)-1%2 == 1 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/getMessageFromMe 10")
		return
	}
	if c.Name == "" || c.Name == "anonymous" {
		// otherwise, prompt user to join to a user
		c.Msg(c, "Okey I got your request but I dont know you.\nPlease Describe your self\n\nHint:)\nname : Specify your name.\n")
		return
	}
	messages, err := s.Service.GetMessageService().GetAllMessages(c.Name)
	if err != nil {
		logrus.WithError(err).Info("GetMessageFromMe error user:", c.Name)
	}
	if len(messages) == 0 {
		c.Msg(c, "You haven't sent a message yet. Now it's time to talk to someone")
		return
	}
	messages = combination(messages, args)
	messageString := ""
	for _, message := range messages {
		messageString += message.ToString()
	}
	c.Msg(c, messageString)

}

// For write to msg last X messages whic is recived to me
func (s *server) getMessageToMe(c *client.Client, args []string) {
	if len(args)-1%2 == 1 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/getMessageTomMe 10")
		return
	}
	if c.Name == "" || c.Name == "anonymous" {
		// otherwise, prompt user to join to a user
		c.Msg(c, "Okey I got your request but I dont know you.\nPlease Describe your self\n\nHint:)\nname : Specify your name.\n")
		return
	}
	messages, err := s.Service.GetMessageService().GetAllMessagesToMe(c.Name)
	if err != nil {
		logrus.WithError(err).Info("GetMessageFromMe error user:", c.Name)
	}
	if len(messages) == 0 {
		c.Msg(c, "You haven't sent a message yet. Now it's time to talk to someone")
		return
	}
	messages = combination(messages, args)
	messageString := ""
	for _, message := range messages {
		messageString += message.ToString()
	}
	c.Msg(c, messageString)

}

// For to write to msg which is last X messages
func (s *server) getLastMassge(c *client.Client, args []string) {
	if len(args)%2 == 1 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/get-last 10")
		return
	}
	if c.Name == "" || c.Name == "anonymous" {
		// otherwise, prompt user to join to a user
		c.Msg(c, "Okey I got your request but I dont know you.\nPlease Describe your self\n\nHint:)\nname : Specify your name.\n")
		return
	}

	_, err := strconv.Atoi(args[1])
	if err != nil {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/get-last 10")
		return
	}
	messages, err := s.Service.GetMessageService().GetLast(c.Name, args[1])
	if err != nil {
		logrus.WithError(err).Info("GetMessageFromMe error user:", c.Name)
	}
	if len(messages) == 0 {
		c.Msg(c, "You haven't sent a message yet. Now it's time to talk to someone")
		return
	}
	messages = combination(messages, args)
	messageString := ""
	for _, message := range messages {
		messageString += message.ToString()
	}
	c.Msg(c, messageString)

}

// For to write to msg which is contains a word
func (s *server) getContains(c *client.Client, args []string) {
	if len(args)%2 == 1 {
		c.Msg(c, "Comand Error: \nCorrect Comamnd Example\n\n/get-last 10")
		return
	}
	if c.Name == "" || c.Name == "anonymous" {
		// otherwise, prompt user to join to a user
		c.Msg(c, "Okey I got your request but I dont know you.\nPlease Describe your self\n\nHint:)\nname : Specify your name.\n")
		return
	}

	messages, err := s.Service.GetMessageService().GetContains(c.Name, args[1])
	if err != nil {
		logrus.WithError(err).Info("GetMessageFromMe error user:", c.Name)
	}
	if len(messages) == 0 {
		c.Msg(c, "You haven't sent a message yet. Now it's time to talk to someone")
		return
	}
	messageString := ""
	for _, message := range messages {
		messageString += message.ToString()
	}
	c.Msg(c, messageString)

}

func combination(messages []model.Message, args []string) []model.Message {
	for i := 1; i < len(args); i += 2 {
		switch args[i] {
		case "||contains":
			var tmpMessages []model.Message
			for _, message := range messages {
				if strings.Contains(message.Text, args[i+1]) {
					tmpMessages = append(tmpMessages, message)
				}
			}
			messages = tmpMessages
		case "||last":
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Println(err)
			}
			if len(messages) <= value {
				break
			}
			var tmpMessages []model.Message
			for index := len(messages) - value; index < len(messages); index++ {
				tmpMessages = append(tmpMessages, messages[index])
			}
			messages = tmpMessages
		}
	}
	return messages
}
