package locserver

import (
	"os"
	"websocket"
	"json"
	l4g "log4go.googlecode.com/hg"
)

var iOpErr = os.NewError("Illegal Operation")
var ider idMaker

// ----------USER------------
// A user is linked 1-1 to a websocket connection
// All of a user's fields can change over its lifetime except writeChan
// Thus, two users are considered equivalent if they have the same writeChan - subject to change
type user struct {
	id        int64          // Globally unique identifier for this user
	tId       int64          // A changing identifier used to tag each message with a transaction id
	Lat, Lng  float64        // Current position of this user, TODO maybe this should be represented as a string?
	mNS, mEW  float64        // Current distances in metres from (lat,lng) (0,0), broken down into North/South and East/West distances
	Name      string         // Arbitrary data (json?) representing this user to the client application
	writeChan chan outPerfer // All messages sent here will be written to the websocket
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
	writeChan := make(chan outPerfer, 32)
	usr := user{id: ider.new(), writeChan: writeChan}
	l4g.Info("User: %d \tConnection Established", usr.id)
	go writeWS(ws, &usr)
	readWS(ws, &usr)
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func readWS(ws *websocket.Conn, usr *user) {
	defer removeOnClose(usr)
	// Accept init message - setName
	buf := make([]byte, 256)
	init, err := unmarshal(usr, buf, ws)
	if err != nil {
		l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		return
	}
	err = processInit(init, usr)
	if err != nil {
		l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		return
	}
	// Accept an endless stream of request messages
	for {
		usr.tId++ // Increase the transaction id for each message
		req, err := unmarshal(usr, buf, ws)
		if err != nil {
			l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
			return
		}
		err = processRequest(req, usr)
		if err != nil {
			l4g.Info("User: %d \tConnection Terminated with %s", usr.id, err.String())
		}
	}
}

// Removes a user from the tree when socket connection is closed
func removeOnClose(usr *user) {
	usr.tId++
	perf := newInPerf(usr.id, usr.tId)
	perf.beginUserProc()
	perf.beginTmSend()
	msgChan <- newCRemove(usr, perf)
}

// Unmarshalls into a *CJsonMsg from the websocket connection returning an error if anything goes wrong
func unmarshal(usr *user, buf []byte, ws *websocket.Conn) (msg *CJsonMsg, err os.Error) {
	n, err := ws.Read(buf)
	if err != nil && err != os.EOF {
		return
	}
	l4g.Info("User: %d \tClient Message: %s", usr.id, string(buf[:n]))
	msg = new(CJsonMsg)
	msg.perf = newInPerf(usr.id, usr.tId)
	msg.perf.beginUserProc()
	err = json.Unmarshal(buf[:n], &msg)
	msg.perf.op = msg.Op
	return
}

// Handle init message - setName
func processInit(init *CJsonMsg, usr *user) (err os.Error) {
	switch init.Op {
	case initOp:
		usr.Name = init.Name
		usr.Lat = init.Lat
		usr.Lng = init.Lng
		usr.mNS, usr.mEW = metresFromOrigin(usr.Lat, usr.Lng)
		init.perf.beginTmSend()
		msgChan <- newCAdd(usr, init.perf)
		msgChan <- newCNearby(usr.Lat, usr.Lng, usr, init.perf)
		return
	}
	return iOpErr
}

// Handle request messages - cRelocate, cNearby
func processRequest(msg *CJsonMsg, usr *user) (err os.Error) {
	switch msg.Op {
	case cNearbyOp:
		lat := msg.Lat
		lng := msg.Lng
		forwardNearby(lat, lng, usr, msg.perf)
		return
	case cMoveOp:
		forwardMove(msg.Lat, msg.Lng, usr, msg.perf)
		return
	}
	return iOpErr
}

// Creates a new cNearby request and sends it to the tree manager
func forwardNearby(lat, lng float64, usr *user, perf *inPerf) {
	nby := newCNearby(lat, lng, usr, perf)
	nby.perf.beginTmSend()
	msgChan <- nby
}

// Creates a new cRelocate request and sends it to the tree manager 
func forwardMove(nLat, nLng float64, usr *user, perf *inPerf) {
	oLat, oLng := usr.Lat, usr.Lng
	// Update usr
	mNS, mEW := metresFromOrigin(nLat, nLng)
	usr.Lat = nLat
	usr.Lng = nLng
	usr.mNS = mNS
	usr.mEW = mEW
	mv := newCMove(oLat, oLng, nLat, nLng, usr, perf)
	mv.perf.beginTmSend()
	msgChan <- mv
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user) {
	writeChan := usr.writeChan
	for {
		v := <-writeChan
		perf := v.getOutPerf()
		perf.beginWSend()
		buf, err := json.MarshalForHTML(v)
		l4g.Info("User: %d \tServer Message: %s", usr.id, string(buf))
		if err != nil {
			l4g.Info("User: %d \tError: %s", usr.id, err.String())
		}
		ws.Write(buf)
		perf.finishAndLog()
	}
}
