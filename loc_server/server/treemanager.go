package locserver

import (
	"log"
	"location_server/quadtree"
	"location_server/msgdef"
	"location_server/profile"
	"location_server/loc_server/user"
)

const (
	nearbyMetresNS = 1000.0
	nearbyMetresEW = 1000.0
)

//
// Single Threaded Tree Manager Code 
//
var msgChan = make(chan clientMsg, 255) // Global Channel for all requests

func TreeManager(minTreeMax int64, trackMovement bool, lg *log.Logger) {
	tree := quadtree.NewQuadTree(maxSouthMetres, maxNorthMetres, maxWestMetres, maxEastMetres, minTreeMax)
	for {
		msg := <-msgChan
		msg.profile.StopAndStart(profile_tmProc)
		switch msg.op {
		case msgdef.CInitLocOp:
			handleInitLoc(&msg, tree, lg)
		case msgdef.CRemoveOp:
			handleRemove(&msg, tree, lg)
		case msgdef.CMoveOp:
			handleMove(&msg, tree, trackMovement, lg)
		case msgdef.CNearbyOp:
			handleNearby(&msg, tree, lg)
		}
		lg.Println(msg.profile.StopAndString())
	}
}

func handleInitLoc(initLoc *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := &initLoc.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLog(usr.Id, "InitLoc Request", mNS, mEW)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, initLocFun(usr))
	tree.Insert(mNS, mEW, usr)
}

func handleRemove(rmv *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := &rmv.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLog(usr.Id, "Remove Request", mNS, mEW)
	deleteUsr(mNS, mEW, usr, tree)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, removeFun(usr))
}

func handleNearby(nby *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := nby.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLog(usr.Id, "Nearby Request", mNS, mEW)
	view := nearbyView(mNS, mEW)
	vs := []*quadtree.View{view}
	tree.Survey(vs, nearbyFun(&usr))
}

func handleMove(mv *clientMsg, tree quadtree.T, trackMovement bool, lg *log.Logger) {
	usr := &mv.usr
	oMNS, oMEW := metresFromOrigin(usr.OLat, usr.OLng)
	nMNS, nMEW := metresFromOrigin(usr.Lat, usr.Lng)
	locLogOldAndNew(usr.Id, "Relocate Request", oMNS, oMEW, nMNS, oMEW)
	deleteUsr(oMNS, oMEW, usr, tree)
	tree.Insert(nMNS, nMEW, usr)
	nView := nearbyView(nMNS, nMEW)
	oView := nearbyView(oMNS, oMEW)
	// Alert out of bounds users
	nvViews := oView.Subtract(nView)
	tree.Survey(nvViews, notVisibleFun(&mv.usr))
	// Alert newly visible users
	vViews := nView.Subtract(oView)
	tree.Survey(vViews, visibleFun(&mv.usr))
	// Alert watching users of the relocation
	if trackMovement {
		movedView := []*quadtree.View{nView.Intersect(oView)}
		tree.Survey(movedView, movedFun(&mv.usr))
	}
}

// Deletes usr from tree at the given coords
func deleteUsr(mNS, mEW float64, usr *user.U, tree quadtree.T) {
	v := quadtree.PointViewP(mNS, mEW)
	pred := func(_, _ float64, e interface{}) bool {
		oUsr := e.(*user.U)
		return usr.Eq(oUsr)
	}
	tree.Delete(v, pred)
}

// Returns a function used for telling usr about each of the other users who are nearby
func nearbyFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Eq(oUsr) {
			broadcastSend(msgdef.SNearbyOp, usr, oUsr)
		}
	}
}

// Returns a function used for alerting users that another user has been added to the system
func initLocFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Eq(oUsr) {
			broadcastSend(msgdef.SAddOp, usr, oUsr)
			broadcastSend(msgdef.SAddOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has been removed from the system
// NB: Relies on the assumption that usr is not currently present in tree
func removeFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		broadcastSend(msgdef.SRemoveOp, usr, oUsr)
	}
}

// Returns a function used for alerting users that another user is going out of range and should be removed
func notVisibleFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		broadcastSend(msgdef.SNotVisibleOp, usr, oUsr)
		broadcastSend(msgdef.SNotVisibleOp, oUsr, usr)
	}
}

// Returns a function used for alerting users that another user has just become visible and should be added
func visibleFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Eq(oUsr) {
			broadcastSend(msgdef.SVisibleOp, usr, oUsr)
			broadcastSend(msgdef.SVisibleOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has changed position and should be updated
func movedFun(usr *user.U) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user.U)
		if !usr.Eq(oUsr) {
			broadcastSend(msgdef.SMovedOp, usr, oUsr)
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

func broadcastSend(op msgdef.ServerOp, usr *user.U, oUsr *user.U) {
	profile := profile.New(usr.Id, usr.TransactionId(), string(op), profile_outTaskNum)
	profile.Start(profile_bSend)
	msg := msgdef.NewServerMsg(op, usr, profile)
	oUsr.WriteMsg(msg)
}

func locLog(id, msgType string, mNS, mEW float64) {
	log.Printf("User: %s \t %s \tmNS: %f mEW: %f", id, msgType, mNS, mEW)
}

func locLogOldAndNew(id, msgType string, oMNS, oMEW, nMNS, nMEW float64) {
	log.Printf("User: %s \t %s \t oMNS: %f oMEW %f nMNS: %f nMEW %f", id, msgType, oMNS, oMEW, nMNS, nMEW)
}
