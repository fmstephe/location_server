package locserver

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/msgwriter"
	"encoding/json"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

type user struct {
	id                   string
	lat, olat, lng, olng float64
	msgWriter            *msgwriter.W
}

func (usr *user) eq(oUsr *user) bool {
	return usr.id == oUsr.id
}

func (usr *user) dup() *user {
	dup := *usr
	return &dup
}

func newUser(ws *websocket.Conn) *user {
	return &user{msgWriter: msgwriter.New(ws)}
}

// A client request
type task struct {
	tId uint
	op  msgdef.ClientOp
	usr *user
}

// NB: The user here is a value, not a pointer
// A copy has been made to avoid race conditions with
// future user updates
func newTask(tId uint, op msgdef.ClientOp, usr *user) *task {
	return &task{tId: tId, op: op, usr: usr.dup()}
}

//  Listen to ws
//  Unmarshall json objects from ws and write to readChan
func WebsocketUser(ws *websocket.Conn) {
	var tId uint
	usr := newUser(ws)
	idMsg := &msgdef.CIdMsg{}
	procReg := processReg(idMsg, usr)
	if err := unmarshalAndProcess(tId, usr.id, ws, idMsg, procReg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	if err := idSet.Add(usr.id, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	logutil.Registered(tId, usr.id)
	defer removeId(&tId, usr)
	tId++
	initLocMsg := msgdef.EmptyCLocMsg()
	procInit := processInitLoc(tId, initLocMsg, usr)
	if err := unmarshalAndProcess(tId, usr.id, ws, initLocMsg, procInit); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	defer removeFromTree(&tId, usr)
	for {
		tId++
		locMsg := msgdef.EmptyCLocMsg()
		procReq := processRequest(tId, locMsg, usr)
		if err := unmarshalAndProcess(tId, usr.id, ws, locMsg, procReq); err != nil {
			usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
			return
		}
	}
}

func unmarshalAndProcess(tId uint, uId string, ws *websocket.Conn, msg interface{}, processFunc func() error) error {
	var data string
	if err := websocket.Message.Receive(ws, &data); err != nil {
		return err
	}
	logutil.Log(tId, uId, data)
	json.Unmarshal([]byte(data), msg)
	return processFunc()
}

func removeId(tId *uint, usr *user) {
	(*tId)++
	logutil.Deregistered(*tId, usr.id)
	idSet.Remove(usr.id)
}

func removeFromTree(tId *uint, usr *user) {
	(*tId)++
	msg := newTask(*tId, msgdef.CRemoveOp, usr)
	forwardMsg(msg)
}

// Handle registration message
// Success will leave usr with initialised Id field
func processReg(idMsg *msgdef.CIdMsg, usr *user) func() error {
	return func() error {
		if err := idMsg.Validate(); err != nil {
			return err
		}
		if idMsg.Op != msgdef.CAddOp {
			return iOpErr
		}
		usr.id = idMsg.Id
		return nil
	}
}

// Handle initial location message
func processInitLoc(tId uint, initMsg *msgdef.CLocMsg, usr *user) func() error {
	return func() error {
		if err := initMsg.Validate(); err != nil {
			return err
		}
		if initMsg.Op != msgdef.CInitLocOp {
			return iOpErr
		}
		usr.olat = initMsg.Lat
		usr.olng = initMsg.Lng
		usr.lat = initMsg.Lat
		usr.lng = initMsg.Lng
		msg := newTask(tId, msgdef.CInitLocOp, usr)
		forwardMsg(msg)
		return nil
	}
}

// Handle request messages - cMove, cNearby
func processRequest(tId uint, locMsg *msgdef.CLocMsg, usr *user) func() error {
	return func() error {
		if err := locMsg.Validate(); err != nil {
			return err
		}
		if locMsg.Op != msgdef.CMoveOp {
			return iOpErr
		}
		usr.olat = usr.lat
		usr.olng = usr.lng
		usr.lat = locMsg.Lat
		usr.lng = locMsg.Lng
		msg := newTask(tId, msgdef.CMoveOp, usr)
		forwardMsg(msg)
		return nil
	}
}

func forwardMsg(msg *task) {
	msgChan <- msg
}
