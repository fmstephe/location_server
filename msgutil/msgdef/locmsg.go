package msgdef

import (
	"errors"
	"math"
)

// Set the initial location of a user
const CInitLocOp = ClientOp("cInitLoc")

// Change the location of the user
const CMoveOp = ClientOp("cMove")

// A structure for unmarshalling lat/lng messages
type CLocMsg struct {
	Op  ClientOp `json:"op"`
	Lat float64  `json:"lat"`
	Lng float64  `json:"lng"`
}

func EmptyCLocMsg() *CLocMsg {
	return &CLocMsg{Lat: math.NaN(), Lng: math.NaN()}
}

func (msg *CLocMsg) Validate() error {
	if msg.Op == "" {
		return errors.New("Missing Op in location message")
	}
	if msg.Op != CInitLocOp && msg.Op != CMoveOp {
		return errors.New("Invalid Op in location message")
	}
	if msg.Lat == math.NaN() || msg.Lng == math.NaN() {
		return errors.New("Lat/Lng position not provided in location message")
	}
	return nil
}

// Indicates that a user has become visible to the receiver
const SVisibleOp = ServerOp("sVisible")

// Indicates that a user has become not visible to the receiver
const SNotVisibleOp = ServerOp("sNotVisible")

// Indicates that a user has moved (and is visible to the receiver)
const SMovedOp = ServerOp("sMoved")

type SLocMsg struct {
	Op  ServerOp `json:"op"`
	Id  string   `json:"id"`
	Lat float64  `json:"lat"`
	Lng float64  `json:"lng"`
}
