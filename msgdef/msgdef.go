package msgdef

import (
	"location_server/profile"
)

// User request operations
type ClientOp string

type ClientMsg struct {
	Op ClientOp
	Msg interface{}
}

func NewClientMsg(op ClientOp, msg interface{}) *ClientMsg {
	return &ClientMsg{Op: op, Msg: msg}
}

// Server reply operations
type ServerOp string


// A server message which contains only a serverOp and a user.CU
type ServerMsg struct {
	Op   ServerOp
	Msg interface{}
}

func NewServerMsg(op ServerOp, msg interface{}) *ServerMsg {
	return &ServerMsg{Op: op, Msg: msg}
}

func NewServerError(msg interface{}) *ServerMsg {
	return &ServerMsg{Op: SErrorOp, Msg: msg}
}

type PServerMsg struct {
	Msg ServerMsg
	Profile *profile.P
}

func NewPServerMsg(op ServerOp, msg interface{}, profile *profile.P) *PServerMsg {
	sm := NewServerMsg(op, msg)
	return &PServerMsg{Msg: *sm, Profile: profile}
}

func NewPServerError(msg interface{}, profile *profile.P) *PServerMsg {
	return NewPServerMsg(SErrorOp, msg, profile)
}

