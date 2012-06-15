package msgdef

import (
	"errors"
	"strings"
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
	if err := validateId(msg.Id); err != nil {
		return err
	}
	return nil
}

func validateId(id string) error {
	if id == "" {
		return errors.New("Id is empty")
	}
	if strings.ContainsAny(id, "<>&'\"") {
		return errors.New("Id contains illegal character(s). May not contain any of <, >, &, ' or \"")
	}
	return nil
}

// Provides a new Id provided by the server
const SIdOp = ServerOp("sId")

// Here it is assumed that the server will not include malicous HTML content
// in user ids
type SIdMsg struct {
	Op ServerOp `json:"op"`
	Id string   `json:"id"`
}
