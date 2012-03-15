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

// Provides a new Id provided by the server
const SIdOp = ServerOp("sId")

type SIdMsg struct {
	Op ServerOp
	Id string
}
