package locserver

import (
	"errors"
	"log"
	"websocket"
	"location_server/msgutil/msgwriter"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/jsonutil"
	"github.com/fmstephe/simpleid"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

type user struct {
	id string
	lat, olat, lng, olng float64
	msgWriter *msgwriter.W
}

func (usr *user) eq(oUsr *user) bool {
	return usr.id == oUsr.id
}

func newUser(ws *websocket.Conn) *user {
	return &user{msgWriter: msgwriter.New(ws)}
}

// A client request
type task struct {
	tId uint
	op   msgdef.ClientOp
	usr  *user
}

func newTask(tId uint, op msgdef.ClientOp, usr *user) *task {
	return &task{tId: tId, op: op, usr: usr}
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func WebsocketUser(ws *websocket.Conn) {
	var tId uint
	var err error
	log.Printf("Connection Established\n")
	usr := newUser(ws)
	idMsg := msgdef.NewCIdMsg()
	if err = unmarshal(ws, idMsg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	if err = processReg(idMsg, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	if err = idSet.Add(usr.id, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	defer removeId(&tId, usr)
	log.Printf("User: %s \tRegistered Successfully\n", usr.id)
	tId++
	initLocMsg := msgdef.NewCLocMsg()
	if err = unmarshal(ws, initLocMsg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	if err = processInitLoc(tId, initLocMsg, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	defer removeFromTree(&tId, usr)
	for {
		tId++
		locMsg := msgdef.NewCLocMsg()
		if err = unmarshal(ws, locMsg); err != nil {
			usr.msgWriter.ErrorAndClose(tId, err.Error())
			return
		}
		if err = processRequest(tId, locMsg, usr); err != nil {
			usr.msgWriter.ErrorAndClose(tId, err.Error())
			return
		}
	}
}

func removeId(tId *uint, usr *user) {
	(*tId)++
	// TODO should probably log this
	idSet.Remove(usr.id)
}

func removeFromTree(tId *uint, usr *user) {
	(*tId)++
	// TODO should probably log this
	msg := newTask(*tId, msgdef.CRemoveOp, usr)
	forwardMsg(msg)
}

// Unmarshals into a *task from the websocket connection returning an error if anything goes wrong
func unmarshal(ws *websocket.Conn, clientMsg *msgdef.ClientMsg) error {
	if err := jsonutil.JSONCodec.Receive(ws, clientMsg); err != nil {
		return err
	}
	//msgLog(usr.id, "Client Message", fmt.Sprintf("%v", clientMsg))
	return nil
}

// Handle registration message
// Function does not return a *task, success will leave usr with initialised Id field
func processReg(clientMsg *msgdef.ClientMsg, usr *user) error {
	idMsg := clientMsg.Msg.(*msgdef.CIdMsg)
	switch clientMsg.Op {
	case msgdef.CAddOp:
		usr.id = idMsg.Id
		return nil
	}
	return iOpErr
}

// Handle initial location message
func processInitLoc(tId uint, clientMsg *msgdef.ClientMsg, usr *user) error {
	initMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CInitLocOp:
		usr.olat = initMsg.Lat
		usr.olng = initMsg.Lng
		usr.lat = initMsg.Lat
		usr.lng = initMsg.Lng
		msg := newTask(tId, msgdef.CInitLocOp, usr)
		forwardMsg(msg)
		return nil
	}
	return iOpErr
}

// Handle request messages - cMove, msg.CNearby
func processRequest(tId uint, clientMsg *msgdef.ClientMsg, usr *user) error {
	locMsg := clientMsg.Msg.(*msgdef.CLocMsg)
	switch clientMsg.Op {
	case msgdef.CNearbyOp:
		msg := newTask(tId, msgdef.CNearbyOp, usr)
		forwardMsg(msg)
		return nil
	case msgdef.CMoveOp:
		usr.olat = usr.lat
		usr.olng = usr.lng
		usr.lat = locMsg.Lat
		usr.lng = locMsg.Lng
		msg := newTask(tId, msgdef.CMoveOp, usr)
		forwardMsg(msg)
		return nil
	}
	return iOpErr
}

func forwardMsg(msg *task) {
	msgChan <- msg
}

func msgLog(id, title, msg string) {
	log.Printf("User: %s\t%s\t%s", id, title, msg)
}
