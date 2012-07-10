package locserver

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"github.com/fmstephe/location_server/logutil"
	"github.com/fmstephe/location_server/msgutil/jsonutil"
	"github.com/fmstephe/location_server/msgutil/msgdef"
	"github.com/fmstephe/location_server/user"
	"github.com/fmstephe/simpleid"
	"math"
)

var iOpErr = errors.New("Illegal Message Op. Operation unrecognised or provided in illegal order.")
var idMap = simpleid.NewIdMap()

// Represents a task for the tree manager.
type task struct {
	tId        uint            // The transaction id for this task
	op         msgdef.ClientOp // The operation to perform for this task
	usr        *user.U         // The state of the user for this task
	olat, olng float64         // The position of the user, if it has changed
}

// Safely creates a new task struct, in particular duplicating usr
func newTask(tId uint, op msgdef.ClientOp, usr *user.U) *task {
	return &task{tId: tId, op: op, usr: usr.Copy(), olat: math.NaN(), olng: math.NaN()}
}

// Safely creates a new task struct, in particular duplicating usr
func newMoveTask(tId uint, op msgdef.ClientOp, usr *user.U, olat, olng float64) *task {
	return &task{tId: tId, op: op, usr: usr.Copy(), olat: olat, olng: olng}
}

// This is the websocket connection handling function
// The following messages are required in this order
// 1: User registration message (user id added to idMap)
// 2: Initial location message 
// 3: Move message
//
// Every incoming message (and subsequent actions performed) are associated with a transaction id
//
// Error handling:
// Any error will result in these actions
// 1: The user will be sent a server-error message
// 2: The connection will be closed
// 3: The user id will be removed from the idMap
// 4: The user will be removed from the treemanager
func HandleLocationService(ws *websocket.Conn) {
	var tId uint
	usr := user.New(ws)
	idMsg := &msgdef.CIdMsg{}
	procReg := processReg(tId, idMsg, usr)
	if err := jsonutil.UnmarshalAndProcess(tId, usr.Id, ws, idMsg, procReg); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	defer removeId(&tId, usr)
	tId++
	initLocMsg := msgdef.EmptyCLocMsg()
	procInit := processInitLoc(tId, initLocMsg, usr)
	if err := jsonutil.UnmarshalAndProcess(tId, usr.Id, ws, initLocMsg, procInit); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	defer removeFromTree(&tId, usr)
	for {
		tId++
		locMsg := msgdef.EmptyCLocMsg()
		procReq := processMove(tId, locMsg, usr)
		if err := jsonutil.UnmarshalAndProcess(tId, usr.Id, ws, locMsg, procReq); err != nil {
			usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
			return
		}
	}
}

// Removes this user's id from idMap and logs the action
func removeId(tId *uint, usr *user.U) {
	(*tId)++
	logutil.Deregistered(*tId, usr.Id)
	idMap.Remove(usr.Id)
}

// Sends a remove message to the tree manager
func removeFromTree(tId *uint, usr *user.U) {
	(*tId)++
	msg := newTask(*tId, msgdef.CRemoveOp, usr)
	forwardMsg(msg)
}

// Handle registration message
// Success will leave usr with initialised Id field
func processReg(tId uint, idMsg *msgdef.CIdMsg, usr *user.U) func() error {
	return func() error {
		if idMsg.Op != msgdef.CAddOp {
			return errors.New("Incorrect op-code for id registration: " + string(idMsg.Op))
		}
		if err := idMsg.Validate(); err != nil {
			return err
		}
		usr.Id = idMsg.Id
		if err := idMap.Add(usr.Id, usr); err != nil {
			return err
		}
		logutil.Registered(tId, usr.Id)
		return nil
	}
}

// Handle initial location message
// Success results in this user's location being updated and an initial location message being sent to the tree manager
func processInitLoc(tId uint, initMsg *msgdef.CLocMsg, usr *user.U) func() error {
	return func() error {
		if err := initMsg.Validate(); err != nil {
			return err
		}
		if initMsg.Op != msgdef.CInitLocOp {
			return iOpErr
		}
		usr.InitLoc(initMsg.Lat, initMsg.Lng)
		msg := newTask(tId, msgdef.CInitLocOp, usr)
		forwardMsg(msg)
		return nil
	}
}

// Handle move message
// Success results in this user's location being updated and a move message beging sent to the tree manager
func processMove(tId uint, locMsg *msgdef.CLocMsg, usr *user.U) func() error {
	return func() error {
		if err := locMsg.Validate(); err != nil {
			return err
		}
		if locMsg.Op != msgdef.CMoveOp {
			return iOpErr
		}
		olat := usr.Lat
		olng := usr.Lng
		usr.Move(locMsg.Lat, locMsg.Lng)
		msg := newMoveTask(tId, msgdef.CMoveOp, usr, olat, olng)
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
