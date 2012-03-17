package msgdef

import (
	"errors"
)

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	Op	ClientOp
	To      string
	Content string
}

func (msg *CMsgMsg) Validate() error {
	if (msg.Op == "") {
		return errors.New("Missing Op in message message")
	}
	if (msg.Op != CMsgOp) {
		return errors.New("Invalid Op in message message")
	}
	if (msg.Content == "") {
		return errors.New("Missing Content in message message")
	}
	return nil
}

// Delivers a message to a user
const SMsgOp = ServerOp("sMsg")

type SMsgMsg struct {
	Op      ServerOp
	From    string
	Content string
}

// Indicates that UserId is not registered on the msg_server
const SNotUserOp = ServerOp("sNotUser")

type SNotUser struct {
	Op     ServerOp
	UserId string
}
