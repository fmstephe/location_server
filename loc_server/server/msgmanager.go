package locserver

import (
	"errors"
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
	fmt.Printf("Connection Established\n")
	go writeWS(ws, usr)
	readWS(ws, usr)
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func readWS(ws *websocket.Conn, usr *user.U) {
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
		usr.NewTransactionId()
		if msg, err := unmarshal(usr, buf, new(msgdef.CLocMsg), processRequest, ws); err != nil {
			fmt.Printf("User: %s \tConnection Terminated with '%s'\n", usr.Id, err.Error())
			return
		} else {
			forwardMsg(msg)
		}
	}
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
	fmt.Printf("User: %s \tClient Message: %s\n", usr.Id, string(buf[:n]))
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
		pMsg := usr.ReceiveMsg()
		profile := pMsg.Profile
		msg := pMsg.Msg
		profile.StopAndStart(profile_wSend)
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
		fmt.Printf("%s\n", profile.StopAndString())
	}
}

func closeWS(ws *websocket.Conn, usr *user.U) {
	if err := ws.Close(); err != nil {
		fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
	}
}
