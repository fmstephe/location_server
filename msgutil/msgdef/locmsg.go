package msgdef

// Set the initial location of a user
const CInitLocOp = ClientOp("cInitLoc")
// Change the location of the user
const CMoveOp = ClientOp("cMove")
// Query for all cNearby users
const CNearbyOp = ClientOp("cNearby")

// A structure for unmarshalling lat/lng messages
type CLocMsg struct {
	Id       string
	Lat, Lng float64
}

func NewCLocMsg() *ClientMsg {
	return &ClientMsg{Msg: &CLocMsg{}}
}

// Indicates that a user has just been added (and is visible to the receiver)
const SAddOp = ServerOp("sAdd")
// Indicates that a user has become visible to the receiver
const SVisibleOp = ServerOp("sVisible")
// Indicates that a user has become not visible to the receiver
const SNotVisibleOp = ServerOp("sNotVisible")
// Indicates that a user has been removed (and was visible to the receiver)
const SRemoveOp = ServerOp("sRemove")
// Indictes that a user is nearby (only sent as response to a nearby request)
const SNearbyOp = ServerOp("sNearby")
// Indicates that a user has moved (and is visible to the receiver)
const SMovedOp = ServerOp("sMoved")
// Indicates that an error has occurred on the server
const SErrorOp = ServerOp("sError")

type SLocMsg struct {
	Id string
	Lat, Lng float64
}
