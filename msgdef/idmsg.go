package msgdef

// Register the id of a user
const CAddOp = ClientOp("cAdd")
// Deregister the id of a user, plus additional cleanup
const CRemoveOp = ClientOp("cRemove")

// A structure for unmarshalling id based messages
type CIdMsg struct {
	Op ClientOp
	Id string
}

func (m *CIdMsg) String() string {
	return string(m.Op)
}

func TestCIdMsg(op ClientOp, id string) *CIdMsg {
	return &CIdMsg{Op: op, Id: id}
}

// Provides a new Id provided by the server
const SIdOp = ServerOp("sId")

// A structure for marshalling new id messages
type SIdMsg struct {
	Op ServerOp
	Id string
}

func NewSIdMsg(id string) *SIdMsg {
	return &SIdMsg{Op: SIdOp, Id: id}
}
