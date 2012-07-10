package locserver

import (
	"fmt"
	"github.com/fmstephe/location_server/logutil"
	"github.com/fmstephe/location_server/msgutil/msgdef"
	"github.com/fmstephe/location_server/quadtree"
	"github.com/fmstephe/location_server/user"
)

// Hardcoded values for the distance within which users are visible to each other
// Should be configurable
const (
	nearbyMetresNS = 1000.0
	nearbyMetresEW = 1000.0
)

// Single channel funnels all messages coming into the tree manager
// As a simple global variable this is a bottleneck (top of the list for performance upgrade)
var taskChan = make(chan *task, 255)

// Starts a goroutine looping listening for messsages on taskChan to process
func StartTreeManager(minTreeMax int64, trackMovement bool) {
	go func() {
		tree := quadtree.NewQuadTree(maxSouthMetres, maxNorthMetres, maxWestMetres, maxEastMetres, minTreeMax)
		for {
			msg := <-taskChan
			switch msg.op {
			case msgdef.CInitLocOp:
				handleInitLoc(msg, tree)
			case msgdef.CRemoveOp:
				handleRemove(msg, tree)
			case msgdef.CMoveOp:
				handleMove(msg, tree, trackMovement)
			}
		}
	}()
}

// Handles initial location tasks
// An initial location message has the following effect
// 1: The user is added to the quadtree at its initial location
// 2: All nearby users to the new user and notified
// 3: Symmetrically the new user is notified of all nearby users
func handleInitLoc(initLoc *task, tree quadtree.T) {
	usr := initLoc.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLog(initLoc.tId, usr.Id, "InitLoc Request", mNS, mEW)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, initLocFun(initLoc.tId, usr))
	tree.Insert(mNS, mEW, usr)
}

// Handles Remove tasks
// A remove task has the following effect
// 1: The user is removed from the quadtree
// 2: All nearby users are notified
func handleRemove(rmv *task, tree quadtree.T) {
	usr := rmv.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLog(rmv.tId, usr.Id, "Remove Request", mNS, mEW)
	deleteUsr(mNS, mEW, usr, tree)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, removeFun(rmv.tId, usr))
}

// Handles move tasks
// A move task has the following effect
// 1: The user is removed from the quadtree at its old location
// 2: The user is inserted into the quadtree at its new location
// 3: All users who could see the user but can't now are notified
// 4: All users who could not see the user but can now are notified
// 5: if (trackMovement) All users who can see the user in both the old and new position are notified
func handleMove(mv *task, tree quadtree.T, trackMovement bool) {
	usr := mv.usr
	oMNS, oMEW := metresFromOrigin(mv.olat, mv.olng)
	nMNS, nMEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLogL(mv.tId, usr.Id, "Relocate Request", oMNS, oMEW, nMNS, oMEW)
	deleteUsr(oMNS, oMEW, usr, tree)
	tree.Insert(nMNS, nMEW, usr)
	nView := nearbyView(nMNS, nMEW)
	oView := nearbyView(oMNS, oMEW)
	// Alert out of bounds users
	nvViews := oView.Subtract(nView)
	tree.Survey(nvViews, notVisibleFun(mv.tId, usr))
	// Alert newly visible users
	vViews := nView.Subtract(oView)
	tree.Survey(vViews, visibleFun(mv.tId, usr))
	// Alert watching users of the relocation
	if trackMovement {
		movedView := []*quadtree.View{nView.Intersect(oView)}
		tree.Survey(movedView, movedFun(mv.tId, usr))
	}
}

// Deletes usr from tree at the given coords
func deleteUsr(mNS, mEW float64, usr *user.U, tree quadtree.T) {
	v := quadtree.PointViewP(mNS, mEW)
	pred := func(_, _ float64, e interface{}) bool {
		oUsr := e.(*user.U)
		return usr.Equiv(oUsr)
	}
	tree.Del(v, pred)
}

// Returns a function used for alerting users that another user has been added to the system
func initLocFun(tId uint, usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Equiv(oUsr) {
			broadcastSend(tId, msgdef.SVisibleOp, usr, oUsr)
			broadcastSend(tId, msgdef.SVisibleOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has been removed from the system
func removeFun(tId uint, usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		broadcastSend(tId, msgdef.SNotVisibleOp, usr, oUsr)
	}
}

// Returns a function used for alerting users that another user has just left the visible range
func notVisibleFun(tId uint, usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		broadcastSend(tId, msgdef.SNotVisibleOp, usr, oUsr)
		broadcastSend(tId, msgdef.SNotVisibleOp, oUsr, usr)
	}
}

// Returns a function used for alerting users that another user has entered the visible range
func visibleFun(tId uint, usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Equiv(oUsr) {
			broadcastSend(tId, msgdef.SVisibleOp, usr, oUsr)
			broadcastSend(tId, msgdef.SVisibleOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user, within visible range, has changed position
func movedFun(tId uint, usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Equiv(oUsr) {
			broadcastSend(tId, msgdef.SMovedOp, usr, oUsr)
		}
	}
}

// Returns a View representing the area considered 'nearby' to the point (mNS,mEW)
func nearbyView(mNS, mEW float64) *quadtree.View {
	sth := mNS - nearbyMetresNS
	nth := mNS + nearbyMetresNS
	wst := mEW - nearbyMetresEW
	est := mEW + nearbyMetresEW
	return quadtree.NewViewP(sth, nth, wst, est)
}

// Sends a message to oUsr informing him/her of a notification involving usr
func broadcastSend(tId uint, op msgdef.ServerOp, usr *user.U, oUsr *user.U) {
	locMsg := msgdef.SLocMsg{Op: op, Id: usr.Id, Lat: usr.Lat, Lng: usr.Lng}
	sMsg := &msgdef.ServerMsg{Msg: locMsg, TId: tId, UId: usr.Id}
	oUsr.MsgWriter.WriteMsg(sMsg)
}

// Logs a task involving only a single location point
func locLog(tId uint, uId, taskDesc string, mNS, mEW float64) {
	logutil.Log(tId, uId, fmt.Sprintf("%s - mNS: %f mEW: %f", taskDesc, mNS, mEW))
}

// Logs a task involving an old and new location point
func locLogL(tId uint, uId, taskDesc string, oMNS, oMEW, nMNS, nMEW float64) {
	logutil.Log(tId, uId, fmt.Sprintf("%s - oMNS: %f oMEW %f nMNS: %f nMEW %f", taskDesc, oMNS, oMEW, nMNS, nMEW))
}
