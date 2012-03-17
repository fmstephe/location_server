package msgdef

import (
	"errors"
	"math"
)

// Set the initial location of a user
const CInitLocOp = ClientOp("cInitLoc")

// Change the location of the user
const CMoveOp = ClientOp("cMove")

// Query for all cNearby users
const CNearbyOp = ClientOp("cNearby")

// A structure for unmarshalling lat/lng messages
type CLocMsg struct {
	Op       ClientOp
	Lat, Lng float64
}

func EmptyCLocMsg() *CLocMsg {
	return &CLocMsg{Lat: math.NaN(), Lng: math.NaN()}
}

func (msg *CLocMsg) Validate() error {
	if msg.Op == "" {
		return errors.New("Missing Op in location message")
	}
	if msg.Op != CInitLocOp && msg.Op != CMoveOp && msg.Op != CNearbyOp {
		return errors.New("Invalid Op in location message")
	}
	if msg.Lat == math.NaN() || msg.Lng == math.NaN() {
		return errors.New("Lat/Lng position not provided in location message")
	}
	return nil
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

type SLocMsg struct {
	Op       ServerOp
	Id       string
	Lat, Lng float64
}
