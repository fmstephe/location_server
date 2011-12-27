package locserver

import (
	"encoding/json"
	"errors"
	"io"
	"websocket"
	"github.com/fmstephe/simpleid"
)

var iOpErr = errors.New("Illegal Operation")
var ider *simpleid.IdMaker

func init() {
	ider = simpleid.New()
}
// ----------USER------------
// A user is linked 1-1 to a websocket connection
// All of a user's fields can change over its lifetime except writeChan
// Thus, two users are considered equivalent if they have the same writeChan - subject to change
type user struct {
	id         int64             // Globally unique identifier for this user
	tId        int64             // A changing identifier used to tag each message with a transaction id
	OLat, OLng float64           // Previous position of this user
	Lat, Lng   float64           // Current position of this user
	Name       string            // Arbitrary data (json?) representing this user to the client application
	writeChan  chan perfProfiler // All messages sent here will be written to the websocket
	// In the future this may contain a reference to the tree it is stored in
}

func (usr *user) eq(oUsr *user) bool {
	if usr == nil || oUsr == nil {
		return false
	}
	return usr.writeChan == oUsr.writeChan
}

//
// Websocket handling per user
//

// Central entry function for a websocket connection
// NB: When we return from this function the websocket will be closed
func WebsocketUser(ws *websocket.Conn) {
	writeChan := make(chan perfProfiler, 32)
	usr := user{id: ider.NewId(), writeChan: writeChan}
	//l4g.Info("User: %d \tConnection Established", usr.id)
	go writeWS(ws, &usr)
	readWS(ws, &usr)
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func readWS(ws *websocket.Conn, usr *user) {
	defer removeOnClose(usr)
	// Accept init message - setName
	buf := make([]byte, 256)
	init, perf, err := unmarshal(usr, buf, ws)
	if err != nil {
		//l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		return
	}
	err = processInit(init, usr, perf)
	if err != nil {
		//l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		return
	}
	// Accept an endless stream of request messages
	for {
		usr.tId++
		req, perf, err := unmarshal(usr, buf, ws)
		if err != nil {
			//l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
			return
		}
		err = processRequest(req, usr, perf)
		if err != nil {
			//l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		}
	}
}

// Removes a user from the tree when socket connection is closed
func removeOnClose(usr *user) {
	usr.tId++
	perf := newPerfProfile(usr.id, usr.tId, string(cRemoveOp), perf_inTaskNum)
	perf.start(perf_userProc)
	forwardMsg(cRemoveOp, usr, perf)
}

// Unmarshalls into a *CJsonMsg from the websocket connection returning an error if anything goes wrong
func unmarshal(usr *user, buf []byte, ws *websocket.Conn) (msg *CJsonMsg, perf *perfProfile, err error) {
	n, err := ws.Read(buf)
	if err != nil && err != io.EOF {
		return
	}
	//l4g.Info("User: %d \tClient Message: %s", usr.id, string(buf[:n]))
	msg = new(CJsonMsg)
	err = json.Unmarshal(buf[:n], &msg)
	perf = newPerfProfile(usr.id, usr.tId, string(msg.Op), perf_inTaskNum)
	perf.start(perf_userProc)
	return
}

// Handle init message - setName
func processInit(init *CJsonMsg, usr *user, perf *perfProfile) (err error) {
	switch init.Op {
	case cAddOp:
		usr.Name = init.Name
		usr.Lat = init.Lat
		usr.Lng = init.Lng
		forwardMsg(cAddOp, usr, perf)
		return
	}
	return iOpErr
}

// Handle request messages - cRelocate, cNearby
func processRequest(msg *CJsonMsg, usr *user, perf *perfProfile) (err error) {
	switch msg.Op {
	case cNearbyOp:
		forwardMsg(cNearbyOp, usr, perf)
		return
	case cMoveOp:
		usr.OLat = usr.Lat
		usr.OLng = usr.Lng
		usr.Lat = msg.Lat
		usr.Lng = msg.Lng
		forwardMsg(cMoveOp, usr, perf)
		return
	}
	return iOpErr
}

func forwardMsg(op clientOp, usr *user, perf *perfProfile) {
	perf.stopAndStart(perf_tmSend)
	msg := newClientMsg(op, usr, perf)
	msgChan <- *msg
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user) {
	writeChan := usr.writeChan
	for {
		msg := <-writeChan
		perf := msg.perfProfile()
		perf.stopAndStart(perf_wSend)
		buf, err := json.MarshalForHTML(msg)
		if err != nil {
			//l4g.Info("User: %d \tError: %s", usr.id, err.String())
		} else {
			//l4g.Info("User: %d \tServer Message: %s", usr.id, string(buf))
			ws.Write(buf)
		}
		//l4g.Info("%s", perf.stopAndString())
	}
}
