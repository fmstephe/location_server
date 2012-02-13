package msgdef

import (
	"location_server/profile"
)

// Set the initial location of a user
const CInitLocOp = ClientOp("cInitLoc")
// Change the location of the user
const CMoveOp = ClientOp("cMove")
// Query for all cNearby users
const CNearbyOp = ClientOp("cNearby")

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

// A structure for unmarshalling lat/lng messages
type CLocMsg struct {
	Op       ClientOp
	Lat, Lng float64
}

func (m *CLocMsg) String() string {
	return string(m.Op)
}

func TestLocMsg(op ClientOp, lat, lng float64) *CLocMsg {
	return &CLocMsg{Op: op, Lat: lat, Lng: lng}
}

// A server message which contains only a serverOp and a user.CU
type ServerMsg struct {
	Op   ServerOp
	Msg interface{}
	profile *profile.P
}

func (m *ServerMsg) Profile() *profile.P {
	return m.profile
}

func NewServerMsg(op ServerOp, msg interface{}, profile *profile.P) *ServerMsg {
	sm := new(ServerMsg)
	sm.Op = op
	sm.Msg = msg
	sm.profile = profile
	return sm
}

func NewServerError(msg interface{}, profile *profile.P) *ServerMsg {
	return NewServerMsg(SErrorOp, msg, profile)
}
