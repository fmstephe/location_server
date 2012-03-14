package msgdef

// User request operations
type ClientOp string

type ClientMsg struct {
	Op  ClientOp
	Msg interface{}
}

func NewClientMsg(op ClientOp, msg interface{}) *ClientMsg {
	return &ClientMsg{Op: op, Msg: msg}
}

// Server reply operations
type ServerOp string

// A server message which contains only a serverOp and a user.CU
type ServerMsg struct {
	Msg interface{}
	TId uint
	UId string
}

// Indicates that an error has occurred on the server
const SErrorOp = ServerOp("sError")

type SErrorMsg struct {
	Op     ServerOp
	ErrMsg string
}

func NewServerError(tId uint, uId string, errMsg string) *ServerMsg {
	msg := &SErrorMsg{Op: SErrorOp, ErrMsg: errMsg}
	return &ServerMsg{Msg: msg, TId: tId, UId: uId}
}
