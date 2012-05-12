package locserver

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/msgwriter"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idSet = simpleid.NewIdMap()

// Identifies a user currently registered with this location service
//
// Valid states (A user moves from one state to the next, and never reverts to a previous state):
// Unregistered: 	msgWriter non-zero
// Registered:		id non-zero, msgWriter non-zero
// Located:		lat/lng fields non-zero, id non-zero, msgWriter non-zero
//
// Since user structs may be part of a goroutines we must take care that only copies are sent.
// NB: It is safe (and necessary) that two copies of the same user reference the same msgWriter
type user struct {
	id                   string
	lat, olat, lng, olng float64
	msgWriter            *msgwriter.W
}

// Indicates whether two user structs indicate the same user
// This is based on the msgWriter pointer as this is the only stable user field
func (usr *user) eq(oUsr *user) bool {
	return usr.msgWriter == oUsr.msgWriter
}

// Duplicates a user for sending as a message
func (usr *user) dup() *user {
	dup := *usr
	return &dup
}

// Creates a new user
func newUser(ws *websocket.Conn) *user {
	return &user{msgWriter: msgwriter.New(ws)}
}

// Represents a task for the tree manager.
type task struct {
	tId uint            // The transaction id for this task
	op  msgdef.ClientOp // The operation to perform for this task
	usr *user           // The state of the user for this task
}

// Safely creates a new task struct, in particular duplicating usr
func newTask(tId uint, op msgdef.ClientOp, usr *user) *task {
	return &task{tId: tId, op: op, usr: usr.dup()}
}

// This is the websocket connection handling function
// The following messages are required in this order
// 1: User registration message (user id added to idSet)
// 2: Initial location message 
// 3: Move message
//
// Every incoming message (and subsequent actions performed) are associated with a transaction id
//
// Error handling:
// Any error will result in these actions
// 1: The user will be sent a server-error message
// 2: The connection will be closed
// 3: The user id will be removed from the idSet
// 4: The user will be removed from the treemanager
func HandleLocationService(ws *websocket.Conn) {
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
		procReq := processMove(tId, locMsg, usr)
		if err := unmarshalAndProcess(tId, usr.id, ws, locMsg, procReq); err != nil {
			usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
			return
		}
	}
}

// Unmarshals a message as a string from the websocket connection
// Unmarshals that string into msg
// Calls processFunc provided for arbitrary handling
func unmarshalAndProcess(tId uint, uId string, ws *websocket.Conn, msg interface{}, processFunc func() error) error {
	var data string
	if err := websocket.Message.Receive(ws, &data); err != nil {
		return err
	}
	logutil.Log(tId, uId, data)
	if err := json.Unmarshal([]byte(data), msg); err != nil {
		return err
	}
	return processFunc()
}

// Removes this user's id from idSet and logs the action
func removeId(tId *uint, usr *user) {
	(*tId)++
	logutil.Deregistered(*tId, usr.id)
	idSet.Remove(usr.id)
}

// Sends a remove message to the tree manager
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
// Success results in this user's location being updated and an initial location message being sent to the tree manager
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

// Handle move message
// Success results in this user's location being updated and  a move message beging sent to the tree manager
func processMove(tId uint, locMsg *msgdef.CLocMsg, usr *user) func() error {
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

// A small function which exists simply to give a level of indirection to this channel send.
// This is clearly a significant bottleneck for the application and in the future this function
// will likely not be a simple channel send.
func forwardMsg(tsk *task) {
	taskChan <- tsk
}
