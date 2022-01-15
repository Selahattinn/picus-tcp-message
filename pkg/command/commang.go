package command

import "github.com/Selahattinn/picus-tcp-message/pkg/client"

// custom type based on int ( for clarity )
type commandID int

const (

	// using iota to generate ever increasing numbers
	CmdName commandID = iota
	CmdJoin
	CmdList
	CmdMsg
	CmdQuit
	CmdHelp
)

// structure for a command
type Command struct {
	ID     commandID
	Client *client.Client
	Args   []string
}
