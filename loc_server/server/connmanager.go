package locserver

import (
	"errors"
	"log"
	"fmt"
	"websocket"
	"encoding/json"
	"location_server/msgdef"
	"location_server/profile"
	"location_server/loc_server/user"
	"github.com/fmstephe/simpleid"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

func jsonMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = json.MarshalForHTML(v)
	return msg, websocket.TextFrame, err
}

func jsonUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return json.Unmarshal(msg, v)
}

var JSONHtml = websocket.Codec{jsonMarshal,jsonUnmarshal}

// A client request
type task struct {
	tId uint
	op   msgdef.ClientOp
	usr  user.U
	profile *profile.P
}

func newTask(tId uint, op msgdef.ClientOp, usr *user.U, profile *profile.P) *task {
	return &task{tId: tId, op: op, usr: *usr, profile: profile}
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
	if _, err := unmarshal(tId, usr, msgdef.NewCIdMsg(), processReg, ws); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	} else {
		log.Printf("User: %s \tRegistered Successfully\n", usr.Id)
	}
	if err := idSet.Add(usr.Id, usr); err != nil {
		writeBackError(tId, usr,"Attempting to insert duplicate user id")
		return
	}
	tId++
	if msg, err := unmarshal(tId, usr, msgdef.NewCLocMsg(), processInitLoc, ws); err != nil {
		writeBackError(tId, usr, err.Error())
		return
	} else {
		forwardMsg(msg)
	}
	defer removeOnClose(&tId, usr)
	for {
		tId++
		if msg, err := unmarshal(tId, usr, msgdef.NewCLocMsg(), processRequest, ws); err != nil {
			writeBackError(tId, usr, err.Error())
			return
		} else {
			forwardMsg(msg)
		}
	}
}

// Log error message and ensure the message is communicated back to the clien (i.e. web-browser etc.)
// The writing go-routine should close the websocket after writing this error message.
func writeBackError(tId uint, usr *user.U, errMsg string) {
	msgLog(usr.Id, "Connection Terminated", errMsg)
	profile := profile.New(profile.ProfileName(tId, usr.Id, errMsg), profile_outTaskNum)
	usr.WriteMsg(msgdef.NewPServerError(errMsg, profile))
	usr.ReceiveClose()
	msgLog(usr.Id, "Close Confirmation Received", "")
}

// Removes a user.U from the tree when socket connection is closed
func removeOnClose(tId *uint, usr *user.U) {
	(*tId)++
	profile := profile.New(profile.ProfileName(*tId, usr.Id, string(msgdef.CRemoveOp)), profile_inTaskNum)
	profile.Start(profile_userProc)
	idSet.Remove(usr.Id)
	msg := newTask(*tId, msgdef.CRemoveOp, usr, profile)
	forwardMsg(msg)
}

// Unmarshals into a *task from the websocket connection returning an error if anything goes wrong
func unmarshal(tId uint, usr *user.U, clientMsg *msgdef.ClientMsg, proc func(uint, *msgdef.ClientMsg, *user.U, *profile.P) (*task, error), ws *websocket.Conn) (task *task, err error) {
	if err = JSONHtml.Receive(ws, clientMsg); err != nil {
		return
	}
	msgLog(usr.Id, "Client Message", fmt.Sprintf("%v", clientMsg))
	profile := profile.New(profile.ProfileName(tId, usr.Id, string(clientMsg.Op)), profile_inTaskNum)
	profile.Start(profile_userProc)
	return proc(tId, clientMsg, usr, profile)
}

// Handle registration message
// Function does not return a *task, success will leave usr with initialised Id field
func processReg(_ uint, clientMsg *msgdef.ClientMsg, usr *user.U, profile *profile.P) (*task, error) {
	idMsg := clientMsg.Msg.(*msgdef.CIdMsg)
	switch clientMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
		return nil, nil
	}
	return nil, iOpErr
}

// Handle initial location message
func processInitLoc(tId uint, clientMsg *msgdef.ClientMsg, usr *user.U, profile *profile.P) (msg *task, err error) {
	initMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CInitLocOp:
		usr.OLat = initMsg.Lat
		usr.OLng = initMsg.Lng
		usr.Lat = initMsg.Lat
		usr.Lng = initMsg.Lng
		msg = newTask(tId, msgdef.CInitLocOp, usr, profile)
		return
	}
	err = iOpErr
	return
}

// Handle request messages - cRelocate, msg.CNearby
func processRequest(tId uint, clientMsg *msgdef.ClientMsg, usr *user.U, profile *profile.P) (msg *task, err error) {
	locMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CNearbyOp:
		msg = newTask(tId, msgdef.CNearbyOp, usr, profile)
		return
	case msgdef.CMoveOp:
		usr.OLat = usr.Lat
		usr.OLng = usr.Lng
		usr.Lat = locMsg.Lat
		usr.Lng = locMsg.Lng
		msg = newTask(tId, msgdef.CMoveOp, usr, profile)
		return
	}
	err = iOpErr
	return
}

func forwardMsg(msg *task) {
	msg.profile.StopAndStart(profile_tmSend)
	msgChan <- msg
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user.U) {
	defer closeWS(ws, usr)
	for {
		pmsg := usr.ReceiveMsg()
		msg := pmsg.Msg
		profile := pmsg.Profile
		profile.StopAndStart(profile_wSend)
		msgLog(usr.Id, "Server Msg", fmt.Sprintf("%v", msg))
		err := JSONHtml.Send(ws, msg)
		log.Printf("%s\n", profile.StopAndString())
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
