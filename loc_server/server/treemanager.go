package locserver

import (
	"fmt"
	"location_server/logutil"
	"location_server/msgutil/msgdef"
	"location_server/quadtree"
)

const (
	nearbyMetresNS = 1000.0
	nearbyMetresEW = 1000.0
)

//
// Single Threaded Tree Manager Code 
//
var msgChan = make(chan *task, 255) // Global Channel for all requests

func TreeManager(minTreeMax int64, trackMovement bool) {
	tree := quadtree.NewQuadTree(maxSouthMetres, maxNorthMetres, maxWestMetres, maxEastMetres, minTreeMax)
	for {
		msg := <-msgChan
		switch msg.op {
		case msgdef.CInitLocOp:
			handleInitLoc(msg, tree)
		case msgdef.CRemoveOp:
			handleRemove(msg, tree)
		case msgdef.CMoveOp:
			handleMove(msg, tree, trackMovement)
		}
	}
}

func handleInitLoc(initLoc *task, tree quadtree.T) {
	usr := initLoc.usr
	mNS, mEW := metresFromOrigin(usr.lat, usr.lng)
	locLog(initLoc.tId, usr.id, "InitLoc Request", mNS, mEW)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, initLocFun(initLoc.tId, usr))
	tree.Insert(mNS, mEW, usr)
}

func handleRemove(rmv *task, tree quadtree.T) {
	usr := rmv.usr
	mNS, mEW := metresFromOrigin(usr.lat, usr.lng)
	locLog(rmv.tId, usr.id, "Remove Request", mNS, mEW)
	deleteUsr(mNS, mEW, usr, tree)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, removeFun(rmv.tId, usr))
}

func handleMove(mv *task, tree quadtree.T, trackMovement bool) {
	usr := mv.usr
	oMNS, oMEW := metresFromOrigin(usr.olat, usr.olng)
	nMNS, nMEW := metresFromOrigin(usr.lat, usr.lng)
	locLogL(mv.tId, usr.id, "Relocate Request", oMNS, oMEW, nMNS, oMEW)
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
func deleteUsr(mNS, mEW float64, usr *user, tree quadtree.T) {
	v := quadtree.PointViewP(mNS, mEW)
	pred := func(_, _ float64, e interface{}) bool {
		oUsr := e.(*user)
		return usr.eq(oUsr)
	}
	tree.Delete(v, pred)
}

// Returns a function used for telling usr about each of the other users who are nearby
func nearbyFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(tId, msgdef.SNearbyOp, usr, oUsr)
		}
	}
}

// Returns a function used for alerting users that another user has been added to the system
func initLocFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(tId, msgdef.SAddOp, usr, oUsr)
			broadcastSend(tId, msgdef.SAddOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has been removed from the system
// NB: Relies on the assumption that usr is not currently present in tree
func removeFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		broadcastSend(tId, msgdef.SRemoveOp, usr, oUsr)
	}
}

// Returns a function used for alerting users that another user is going out of range and should be removed
func notVisibleFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		broadcastSend(tId, msgdef.SNotVisibleOp, usr, oUsr)
		broadcastSend(tId, msgdef.SNotVisibleOp, oUsr, usr)
	}
}

// Returns a function used for alerting users that another user has just become visible and should be added
func visibleFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(tId, msgdef.SVisibleOp, usr, oUsr)
			broadcastSend(tId, msgdef.SVisibleOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has changed position and should be updated
func movedFun(tId uint, usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
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

func broadcastSend(tId uint, op msgdef.ServerOp, usr *user, oUsr *user) {
	locMsg := msgdef.SLocMsg{Op: op, Id: usr.id, Lat: usr.lat, Lng: usr.lng}
	sMsg := &msgdef.ServerMsg{Msg: locMsg, TId: tId, UId: usr.id}
	oUsr.msgWriter.WriteMsg(sMsg)
}

func locLog(tId uint, uId, msgType string, mNS, mEW float64) {
	logutil.Log(tId, uId, fmt.Sprintf("%s - mNS: %f mEW: %f", msgType, mNS, mEW))
}

func locLogL(tId uint, uId, msgType string, oMNS, oMEW, nMNS, nMEW float64) {
	logutil.Log(tId, uId, fmt.Sprintf("%s - oMNS: %f oMEW %f nMNS: %f nMEW %f", msgType, oMNS, oMEW, nMNS, nMEW))
}
