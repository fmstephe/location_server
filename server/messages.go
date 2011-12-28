package locserver

// ----------CLIENT MESSAGES------------
// User request operations NB: a user is synonymous with a websocket connection
type clientOp string

// Set the name of the user and store user location in the tree
const cAddOp = clientOp("cAdd")
// Remove user from the tree - result of closed connection
const cRemoveOp = clientOp("cRemove")
// Change the location of the user
const cMoveOp = clientOp("cMove")
// Query for all cNearby users
const cNearbyOp = clientOp("cNearby")

// A structure for reading messages off a websocket
type CJsonMsg struct {
	Op       clientOp
	Lat, Lng float64
	Id     string
}

// Usually CJsonMsg structs are built by unmarshal - but this is hard for external test packages
func TestMsg(op clientOp, lat, lng float64, id string) *CJsonMsg {
	return &CJsonMsg{Op: op, Lat: lat, Lng: lng, Id: id}
}

// A client request
type clientMsg struct {
	op   clientOp
	usr  user
	perf perfProfile
}

func newClientMsg(op clientOp, usr *user, perf *perfProfile) *clientMsg {
	m := new(clientMsg)
	m.op = op
	m.usr = *usr
	m.perf = *perf
	return m
}

// ----------SERVER MESSAGES------------
// Server reply operations
type serverOp string

// Indicates that a user has just been added (and is visible to the receiver)
const sAddOp = serverOp("sAdd")
// Indicates that a user has become visible to the receiver
const sVisibleOp = serverOp("sVisible")
// Indicates that a user has become not visible to the receiver
const sNotVisibleOp = serverOp("sNotVisible")
// Indicates that a user has been removed (and was visible to the receiver)
const sRemoveOp = serverOp("sRemove")
// Indictes that a user is nearby (only sent as response to a nearby request)
const sNearbyOp = serverOp("sNearby")
// Indicates that a user has moved (and is visible to the receiver)
const sMovedOp = serverOp("sMoved")

// A server message which contains only a serverOp and a user
type serverMsg struct {
	Op   serverOp
	Usr  user
	perf perfProfile
}

func (msg *serverMsg) perfProfile() *perfProfile {
	return &msg.perf
}

func newServerMsg(op serverOp, usr *user, perf *perfProfile) *serverMsg {
	m := new(serverMsg)
	m.Op = op
	m.Usr = *usr
	m.perf = *perf
	return m
}
