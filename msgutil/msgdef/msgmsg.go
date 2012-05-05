package msgdef

import (
	"errors"
)

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	Op      ClientOp `json:"op"`
	To      string   `json:"to"`
	Content string   `json:"content"`
}

func (msg *CMsgMsg) Validate() error {
	if msg.Op == "" {
		return errors.New("Missing Op in message message")
	}
	if msg.Op != CMsgOp {
		return errors.New("Invalid Op in message message")
	}
	if msg.Content == "" {
		return errors.New("Missing Content in message message")
	}
	return nil
}

// Delivers a message to a user
const SMsgOp = ServerOp("sMsg")

type SMsgMsg struct {
	Op      ServerOp `json:"op"`
	From    string   `json:"from"`
	Content string   `json:"content"`
}

// Indicates that UserId is not registered on the msg_server
const SNotUserOp = ServerOp("sNotUser")

type SNotUser struct {
	Op     ServerOp `json:"op"`
	UserId string   `json:"user-id"`
}
