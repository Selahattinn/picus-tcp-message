package model

import "strconv"

type Message struct {
	ID   int64
	From string
	To   string
	Text string
}

func (m Message) ToString() string {
	message := ""
	message = "ID: " + strconv.FormatInt(m.ID, 10) + "\n\tFrom: " + m.From + "\n\tTo: " + m.To + "\n\tmessage: " + m.Text + "\n"
	return message
}
