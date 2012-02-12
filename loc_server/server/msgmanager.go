package locserver

import (
	"errors"
	"log"
	"fmt"
	"io"
	"websocket"
	"encoding/json"
	"location_server/msgdef"
	"location_server/profile"
	"location_server/loc_server/user"
	"github.com/fmstephe/simpleid"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

// A client request
type clientMsg struct {
	op   msgdef.ClientOp
	usr  user.U
	profile profile.P
}

func newClientMsg(op msgdef.ClientOp, usr *user.U, profile profile.P) *clientMsg {
	return &clientMsg{op: op, usr: *usr, profile: profile}
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
	buf := make([]byte, 256)
	if _, err := unmarshal(usr, buf, new(msgdef.CIdMsg), processReg, ws); err != nil {
		writeBackError(usr, err.Error())
		return
	} else {
		log.Printf("User: %s \tRegistered Successfully\n", usr.Id)
	}
	if err := idSet.Add(usr.Id, usr); err != nil {
		writeBackError(usr,"Attempting to insert duplicate user id")
		return
	}
	if msg, err := unmarshal(usr, buf, new(msgdef.CLocMsg), processInitLoc, ws); err != nil {
		writeBackError(usr, err.Error())
		return
	} else {
		forwardMsg(msg)
	}
	defer removeOnClose(usr)
	for {
		usr.NewTransactionId()
		if msg, err := unmarshal(usr, buf, new(msgdef.CLocMsg), processRequest, ws); err != nil {
			writeBackError(usr, err.Error())
			return
		} else {
			forwardMsg(msg)
		}
	}
}

// Log error message and ensure the message is communicated back to the clien (i.e. web-browser etc.)
// The writing go-routine should close the websocket after writing this error message.
func writeBackError(usr *user.U, errMsg string) {
	msgLog(usr.Id, "Connection Terminated", errMsg)
	profile := profile.New(usr.Id, usr.TransactionId(), errMsg, profile_outTaskNum)
	usr.WriteMsg(msgdef.NewServerError(errMsg, profile))
	usr.ReceiveClose()
}

// Removes a user.U from the tree when socket connection is closed
func removeOnClose(usr *user.U) {
	tId := usr.NewTransactionId()
	profile := profile.New(usr.Id, tId, string(msgdef.CRemoveOp), profile_inTaskNum)
	profile.Start(profile_userProc)
	idSet.Remove(usr.Id)
	msg := newClientMsg(msgdef.CRemoveOp, usr, profile)
	forwardMsg(msg)
}

// Unmarshals into a *CMsg from the websocket connection returning an error if anything goes wrong
func unmarshal(usr *user.U, buf []byte, rMsg fmt.Stringer, proc func(fmt.Stringer, *user.U, profile.P) (*clientMsg, error), ws *websocket.Conn) (*clientMsg, error) {
	n, err := ws.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	msgLog(usr.Id, "Client Message", string(buf[:n]))
	if err = json.Unmarshal(buf[:n], rMsg); err != nil {
		return nil, err
	}
	profile := profile.New(usr.Id, usr.TransactionId(), rMsg.String(), profile_inTaskNum)
	profile.Start(profile_userProc)
	return proc(rMsg, usr, profile)
}

// Handle registration message
// Function does not return a *clientMsg, success will leave usr with initialised Id field
func processReg(rMsg fmt.Stringer, usr *user.U, profile profile.P) (*clientMsg, error) {
	idMsg := rMsg.(*msgdef.CIdMsg)
	switch idMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
		return nil, nil
	}
	return nil, iOpErr
}

// Handle initial location message
func processInitLoc(rMsg fmt.Stringer, usr *user.U, profile profile.P) (msg *clientMsg, err error) {
	initMsg := rMsg.(*msgdef.CLocMsg)
	switch initMsg.Op {
	case msgdef.CInitLocOp:
		usr.OLat = initMsg.Lat
		usr.OLng = initMsg.Lng
		usr.Lat = initMsg.Lat
		usr.Lng = initMsg.Lng
		msg = newClientMsg(msgdef.CInitLocOp, usr, profile)
		return
	}
	err = iOpErr
	return
}

// Handle request messages - cRelocate, msg.CNearby
func processRequest(rMsg fmt.Stringer, usr *user.U, profile profile.P) (msg *clientMsg, err error) {
	locMsg := rMsg.(*msgdef.CLocMsg)
	switch locMsg.Op {
	case msgdef.CNearbyOp:
		msg = newClientMsg(msgdef.CNearbyOp, usr, profile)
		return
	case msgdef.CMoveOp:
		usr.OLat = usr.Lat
		usr.OLng = usr.Lng
		usr.Lat = locMsg.Lat
		usr.Lng = locMsg.Lng
		msg = newClientMsg(msgdef.CMoveOp, usr, profile)
		return
	}
	err = iOpErr
	return
}

func forwardMsg(msg *clientMsg) {
	msg.profile.StopAndStart(profile_tmSend)
	msgChan <- *msg
}

//  Listen to writeChan
//  Marshal values from writeChan and write to ws
func writeWS(ws *websocket.Conn, usr *user.U) {
	defer closeWS(ws, usr)
	for {
		msg := usr.ReceiveMsg()
		profile := msg.Profile()
		profile.StopAndStart(profile_wSend)
		buf, err := json.MarshalForHTML(msg)
		if err != nil {
			msgLog(usr.Id, "Error", err.Error())
			return
		}
		msgLog(usr.Id, "Server Message", err.Error())
		if _, err = ws.Write(buf); err != nil {
			msgLog(usr.Id, "Error", err.Error())
			return
		}
		log.Printf("%s\n", profile.StopAndString())
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

func msgLog(id string, title, msg string) {
	log.Printf("User: %s\t%s\t%s", id, title, msg)
}
