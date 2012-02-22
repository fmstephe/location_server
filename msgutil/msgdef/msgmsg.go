package msgdef

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	To      string
	Content string
}

func NewCMsgMsg() *ClientMsg {
	return &ClientMsg{Msg: &CMsgMsg{}}
}

// Delivers a message to a user
const SMsgOp = ServerOp("sMsg")

type SMsgMsg struct {
	From    string
	Content string
}
