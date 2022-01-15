package client

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Selahattinn/picus-tcp-message/pkg/command"
	"github.com/Selahattinn/picus-tcp-message/pkg/crypto"
)

// structure of a client i.e a user ( a new connection, will have this structure )
type Client struct {

	// client connection details
	Conn net.Conn

	// identifier for the client :
	// client will be known on the server by this name
	Name string

	// identifier for the (other) client :
	// the other client (person), this client is talking to currently
	Contact string

	// commands to facilitate chat system
	Commands chan<- command.Command

	// private
	Private *rsa.PrivateKey

	// public
	Public rsa.PublicKey
}

// function to read input
func (c *Client) ReadInput() {

	// continuously...
	for {

		// read user input
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			// abort if an error occurs
			return
		}

		// process input, to parse commands
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		// update client command for desired command
		switch cmd {

		case "/name":
			// specify your name
			c.Commands <- command.Command{
				ID:     command.CmdName,
				Client: c,
				Args:   args,
			}
		case "/join":
			// connect to another user :
			// to be able to chat with him/her
			c.Commands <- command.Command{
				ID:     command.CmdJoin,
				Client: c,
				Args:   args,
			}
		case "/list":
			// display all the available users on the server :
			// these are ones you ( a client ) can join and chat to
			c.Commands <- command.Command{
				ID:     command.CmdList,
				Client: c,
			}
		case "/msg":
			// send a message to the user ( another client ) you have joined
			c.Commands <- command.Command{
				ID:     command.CmdMsg,
				Client: c,
				Args:   args,
			}
		case "/quit":
			// exit the chat system
			c.Commands <- command.Command{
				ID:     command.CmdQuit,
				Client: c,
			}
		case "/help":
			// return command list
			c.Commands <- command.Command{
				ID:     command.CmdHelp,
				Client: c,
			}
			// for any other command
		default:
			c.Err(fmt.Errorf("unknown command: %s", cmd))
			c.Msg(c, "* use '/help' to   command to send message to selected user")
		}
	}
}

// writes an error message current client
func (c *Client) Err(err error) {

	_, e := c.Conn.Write([]byte("err: " + err.Error() + "\n"))
	if e != nil {
		log.Fatalln("unable to write to connection", e)
	}
}

// writes a message to specified client
func (c *Client) Msg(x *Client, msg string) {

	// if contacting other client
	if c.Private != x.Private {

		dMsg := crypto.Decrypt(msg, *x.Private)

		// write message to client
		_, e := x.Conn.Write([]byte("> " + dMsg + "\n"))
		if e != nil {
			log.Fatalln("unable to write over client connection")
		}

	} else {
		// write message to client
		_, e := x.Conn.Write([]byte("> " + msg + "\n"))
		if e != nil {
			log.Fatalln("unable to write over client connection")
		}
	}

}
