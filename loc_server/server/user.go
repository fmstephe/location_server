package locserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fmstephe/simpleid"
	"io"
	"location_server/msgdef"
	"location_server/profile"
	"websocket"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

type user struct {
	Id         string          // Unique identifier for this user
	tId        int64           // A changing identifier used to tag each message with a transaction id
	OLat, OLng float64         // Previous position of this user
	Lat, Lng   float64         // Current position of this user
	writeChan  chan *serverMsg // All messages sent here will be written to the websocket
	// In the future this may contain a reference to the tree it is stored in
}

// A client request
type clientMsg struct {
	op   msgdef.ClientOp
	usr  user
	perf profile.P
}

func newClientMsg(op msgdef.ClientOp, usr *user, perf profile.P) *clientMsg {
	return &clientMsg{op: op, usr: *usr, perf: perf}
}

// A server message which contains only a serverOp and a user
type serverMsg struct {
	Op   msgdef.ServerOp
	Usr  user
	perf profile.P
}

func newServerMsg(op msgdef.ServerOp, usr *user, perf profile.P) *serverMsg {
	m := new(serverMsg)
	m.Op = op
	m.Usr = *usr
	m.perf = perf
	return m
}

func (usr *user) eq(oUsr *user) bool {
	if usr == nil || oUsr == nil {
		return false
	}
	return usr.Id == oUsr.Id
}

// Central entry function for a websocket connection
// NB: When we return from this function the websocket will be closed
func WebsocketUser(ws *websocket.Conn) {
	writeChan := make(chan *serverMsg, 32)
	usr := user{writeChan: writeChan}
	fmt.Printf("Connection Established\n")
	go writeWS(ws, &usr)
	readWS(ws, &usr)
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func readWS(ws *websocket.Conn, usr *user) {
	buf := make([]byte, 256)
	if _, err := unmarshal(usr, buf, new(msgdef.CIdMsg), processReg, ws); err != nil {
		fmt.Printf("User: %s \tConnection Terminated with '%s'\n", usr.Id, err.Error())
		return
	} else {
		fmt.Printf("User: %s \tRegistered Successfully\n", usr.Id)
	}
	if err := idSet.Add(usr.Id, usr); err != nil {
		return
	}
	if msg, err := unmarshal(usr, buf, new(msgdef.CLocMsg), processInitLoc, ws); err != nil {
		fmt.Printf("User: %s \tConnection Terminated with '%s'\n", usr.Id, err.Error())
		return
	} else {
		forwardMsg(msg)
	}
	defer removeOnClose(usr)
	for {
		usr.tId++
		if msg, err := unmarshal(usr, buf, new(msgdef.CLocMsg), processRequest, ws); err != nil {
			fmt.Printf("User: %s \tConnection Terminated with '%s'\n", usr.Id, err.Error())
			return
		} else {
			forwardMsg(msg)
		}
	}
}

// Removes a user from the tree when socket connection is closed
func removeOnClose(usr *user) {
	usr.tId++
	perf := profile.New(usr.Id, usr.tId, string(msgdef.CRemoveOp), profile_inTaskNum)
	perf.Start(profile_userProc)
	idSet.Remove(usr.Id)
	msg := newClientMsg(msgdef.CRemoveOp, usr, perf)
	forwardMsg(msg)
}

// Unmarshals into a *CMsg from the websocket connection returning an error if anything goes wrong
func unmarshal(usr *user, buf []byte, rMsg fmt.Stringer, proc func(fmt.Stringer, *user, profile.P) (*clientMsg, error), ws *websocket.Conn) (*clientMsg, error) {
	n, err := ws.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	fmt.Printf("User: %s \tClient Message: %s\n", usr.Id, string(buf[:n]))
	if err = json.Unmarshal(buf[:n], rMsg); err != nil {
		return nil, err
	}
	perf := profile.New(usr.Id, usr.tId, rMsg.String(), profile_inTaskNum)
	perf.Start(profile_userProc)
	return proc(rMsg, usr, perf)
}

// Handle registration message
// Function does not return a *clientMsg, success will leave usr with initialised Id field
func processReg(rMsg fmt.Stringer, usr *user, perf profile.P) (*clientMsg, error) {
	idMsg := rMsg.(*msgdef.CIdMsg)
	switch idMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
		return nil, nil
	}
	return nil, iOpErr
}

// Handle initial location message
func processInitLoc(rMsg fmt.Stringer, usr *user, perf profile.P) (msg *clientMsg, err error) {
	initMsg := rMsg.(*msgdef.CLocMsg)
	switch initMsg.Op {
	case msgdef.CInitLocOp:
		usr.OLat = initMsg.Lat
		usr.OLng = initMsg.Lng
		usr.Lat = initMsg.Lat
		usr.Lng = initMsg.Lng
		msg = newClientMsg(msgdef.CInitLocOp, usr, perf)
		return
	}
	err = iOpErr
	return
}

// Handle request messages - cRelocate, msg.CNearby
func processRequest(rMsg fmt.Stringer, usr *user, perf profile.P) (msg *clientMsg, err error) {
	locMsg := rMsg.(*msgdef.CLocMsg)
	switch locMsg.Op {
	case msgdef.CNearbyOp:
		msg = newClientMsg(msgdef.CNearbyOp, usr, perf)
		return
	case msgdef.CMoveOp:
		usr.OLat = usr.Lat
		usr.OLng = usr.Lng
		usr.Lat = locMsg.Lat
		usr.Lng = locMsg.Lng
		msg = newClientMsg(msgdef.CMoveOp, usr, perf)
		return
	}
	err = iOpErr
	return
}

func forwardMsg(msg *clientMsg) {
	msg.perf.StopAndStart(profile_tmSend)
	msgChan <- *msg
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user) {
	writeChan := usr.writeChan
	defer closeWS(ws, usr)
	for {
		msg := <-writeChan
		perf := msg.perf
		perf.StopAndStart(profile_wSend)
		buf, err := json.MarshalForHTML(msg)
		if err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
			return
		}
		fmt.Printf("User: %s \tServer Message: %s\n", usr.Id, string(buf))
		if _, err = ws.Write(buf); err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
			return
		}
		fmt.Printf("%s\n", perf.StopAndString())
	}
}

func closeWS(ws *websocket.Conn, usr *user) {
	if err := ws.Close(); err != nil {
		fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
	}
}
