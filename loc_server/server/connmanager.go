package locserver

import (
	"errors"
	"log"
	"fmt"
	"websocket"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/jsonutil"
	"location_server/loc_server/user"
	"github.com/fmstephe/simpleid"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

type locUsr struct {
	lat, olat, lng, olng float64
	usr *user.U
}

// A client request
type task struct {
	tId uint
	op   msgdef.ClientOp
	usr  user.U
}

func newTask(tId uint, op msgdef.ClientOp, usr *user.U) *task {
	return &task{tId: tId, op: op, usr: *usr}
}

// Central entry function for a websocket connection
// NB: When we return from this function the websocket will be closed
func WebsocketUser(ws *websocket.Conn) {
	usr := user.New()
	log.Printf("Connection Established\n")
	go writeWS(ws, usr)
	readWS(ws, usr)
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func readWS(ws *websocket.Conn, usr *user.U) {
	var tId uint
	var err error
	idMsg := msgdef.NewCIdMsg()
	if err = unmarshal(tId, usr, idMsg, ws); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	}
	if err = processReg(idMsg, usr, ); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	}
	if err = idSet.Add(usr.Id, usr); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	}
	log.Printf("User: %s \tRegistered Successfully\n", usr.Id)
	tId++
	initLocMsg := msgdef.NewCLocMsg()
	if err = unmarshal(tId, usr, initLocMsg, ws); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	}
	if err = processInitLoc(tId, initLocMsg, usr, ); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	}
	defer removeOnClose(&tId, usr)
	for {
		tId++
		locMsg := msgdef.NewCLocMsg()
		if err = unmarshal(tId, usr, locMsg, ws); err != nil {
			writeBackError(tId, usr, err.Error())
			return
		}
		if err = processRequest(tId, locMsg, usr, ); err != nil {
			writeBackError(tId, usr, err.Error())
		}
	}
}

// Log error message and ensure the message is communicated back to the client (i.e. web-browser etc.)
// The writing go-routine should close the websocket after writing this error message.
func writeBackError(tId uint, usr *user.U, errMsg string) {
	msgLog(usr.Id, "Connection Terminated", errMsg)
	usr.WriteMsg(msgdef.NewServerError(errMsg))
	usr.ReceiveClose()
	msgLog(usr.Id, "Close Confirmation Received", "")
}

// Removes a user.U from the tree when socket connection is closed
func removeOnClose(tId *uint, usr *user.U) {
	(*tId)++
	idSet.Remove(usr.Id)
	msg := newTask(*tId, msgdef.CRemoveOp, usr)
	forwardMsg(msg)
}

// Unmarshals into a *task from the websocket connection returning an error if anything goes wrong
func unmarshal(tId uint, usr *user.U, clientMsg *msgdef.ClientMsg, ws *websocket.Conn) error {
	if err := jsonutil.JSONCodec.Receive(ws, clientMsg); err != nil {
		return err
	}
	msgLog(usr.Id, "Client Message", fmt.Sprintf("%v", clientMsg))
	return nil
}

// Handle registration message
// Function does not return a *task, success will leave usr with initialised Id field
func processReg(clientMsg *msgdef.ClientMsg, usr *user.U) error {
	idMsg := clientMsg.Msg.(*msgdef.CIdMsg)
	switch clientMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
		return nil
	}
	return iOpErr
}

// Handle initial location message
func processInitLoc(tId uint, clientMsg *msgdef.ClientMsg, usr *user.U) error {
	initMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CInitLocOp:
		usr.OLat = initMsg.Lat
		usr.OLng = initMsg.Lng
		usr.Lat = initMsg.Lat
		usr.Lng = initMsg.Lng
		msg := newTask(tId, msgdef.CInitLocOp, usr)
		forwardMsg(msg)
		return nil
	}
	return iOpErr
}

// Handle request messages - cRelocate, msg.CNearby
func processRequest(tId uint, clientMsg *msgdef.ClientMsg, usr *user.U) error {
	locMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CNearbyOp:
		msg := newTask(tId, msgdef.CNearbyOp, usr)
		forwardMsg(msg)
		return nil
	case msgdef.CMoveOp:
		usr.OLat = usr.Lat
		usr.OLng = usr.Lng
		usr.Lat = locMsg.Lat
		usr.Lng = locMsg.Lng
		msg := newTask(tId, msgdef.CMoveOp, usr)
		forwardMsg(msg)
		return nil
	}
	return iOpErr
}

func forwardMsg(msg *task) {
	msgChan <- msg
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user.U) {
	defer closeWS(ws, usr)
	for {
		msg := usr.ReceiveMsg()
		msgLog(usr.Id, "Server Msg", fmt.Sprintf("%v", msg))
		err := jsonutil.JSONCodec.Send(ws, msg)
		if err != nil {
			msgLog(usr.Id, "Error", err.Error())
			return
		}
		if msg.Op == msgdef.SErrorOp {
			closeWS(ws, usr)
			return
		}
	}
}

func closeWS(ws *websocket.Conn, usr *user.U) {
	if err := ws.Close(); err != nil {
		msgLog(usr.Id, "Error", err.Error())
	}
	usr.WriteClose()
}

func msgLog(id, title, msg string) {
	log.Printf("User: %s\t%s\t%s", id, title, msg)
}
