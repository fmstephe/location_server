package msgdef

// Sends a message to another user
const CMsgOp = ClientOp("cMsg")

type CMsgMsg struct {
	Op	ClientOp
	To      string
	Content string
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
