package msgdef

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	Op  ClientOp
	For string
	Msg string
}

// Delivers a message to a user
const SMsgOp = ServerOp("sMsg")

type SMsgMsg struct {
	Op   ServerOp
	From string
	Msg  string
}
