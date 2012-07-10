package user

import (
	"code.google.com/p/go.net/websocket"
	"github.com/fmstephe/location_server/msgutil/msgwriter"
)

// Identifies a user identifed with a lat/lng location currently registered with this service
//
// Valid states (A user moves from one state to the next, and never reverts to a previous state):
// Unregistered:        MsgWriter non-zero
// Registered:          id non-zero, MsgWriter non-zero
// Located:             lat/lng fields non-zero, id non-zero, MsgWriter non-zero
//
// U is not thread safe. When sharing between goroutines care must be taken to only send copies.
// NB: It is safe (and necessary) that two copies of the same user reference the same MsgWriter
type U struct {
	Id        string
	Lat, Lng  float64
	MsgWriter *msgwriter.W
}

// Initialises the location of a user
func (usr *U) InitLoc(lat, lng float64) {
	usr.Lat = lat
	usr.Lng = lng
}

// Moves the user to a new location. The previous location is saved in the OLat/OLng fields.
func (usr *U) Move(lat, lng float64) {
	usr.Lat = lat
	usr.Lng = lng
}

// Indicates whether two user structs indicate the same user
// This is based on the MsgWriter pointer as this is the only stable user field
func (usr *U) Equiv(oUsr *U) bool {
	return usr.MsgWriter == oUsr.MsgWriter
}

// Duplicates a user for sending as a message
func (usr *U) Copy() *U {
	dup := *usr
	return &dup
}

// Creates a new user
func New(ws *websocket.Conn) *U {
	return &U{MsgWriter: msgwriter.New(ws)}
}
