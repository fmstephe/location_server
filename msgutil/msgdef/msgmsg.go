package msgdef

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	To      string
	Id      string
	sends   uint
	Content string
}

func NewCMsgMsg() *ClientMsg {
	return &ClientMsg{Msg: &CMsgMsg{}}
}

// Delivers a message to a user
const SMsgOp = ServerOp("sMsg")

type SMsgMsg struct {
	From    string
	Id      string
	sends   uint
	Content string
}

// Indicates that the user 
const SNotUserOp = ServerOp("sNotUser")

type SNotUser struct {
	UserId string
}

// Acknowledges a message was received
const SAckOp = ServerOp("sAck")

type SAckMsg struct {
	From string
	Id string
}
