package msgdef

// Register the id of a user
const CAddOp = ClientOp("cAdd")

// Deregister the id of a user, plus additional cleanup
const CRemoveOp = ClientOp("cRemove")

// A structure for unmarshalling id based messages
type CIdMsg struct {
	Id string
}

func NewCIdMsg() *ClientMsg {
	return &ClientMsg{Msg: &CIdMsg{}}
}

// Provides a new Id provided by the server
const SIdOp = ServerOp("sId")

type SIdMsg struct {
	Op ServerOp
	Id string
}
