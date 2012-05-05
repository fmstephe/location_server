package msgdef

import (
	"errors"
)

// Register the id of a user
const CAddOp = ClientOp("cAdd")

// Deregister the id of a user, plus additional cleanup
const CRemoveOp = ClientOp("cRemove")

// A structure for unmarshalling id based messages
type CIdMsg struct {
	Op ClientOp `json:"op"`
	Id string   `json:"id"`
}

func (msg *CIdMsg) Validate() error {
	if msg.Op == "" {
		return errors.New("Missing Op in id message")
	}
	if msg.Op != CAddOp && msg.Op != CRemoveOp {
		return errors.New("Invalid Op in id message")
	}
	if msg.Id == "" {
		return errors.New("Missing Id in id message")
	}
	return nil
}

// Provides a new Id provided by the server
const SIdOp = ServerOp("sId")

type SIdMsg struct {
	Op ServerOp `json:"op"`
	Id string   `json:"id"`
}
