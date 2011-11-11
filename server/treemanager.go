package locserver

import (
	"location_server/quadtree"
	"log"
)

const (
	nearbyMetresNS = 1000.0
	nearbyMetresEW = 1000.0
)

//
// Single Threaded Tree Manager Code 
//
var msgChan = make(chan clientMsg, 255) // Global Channel for all requests

func TreeManager(minTreeMax int64, lg *log.Logger) {
	tree := quadtree.NewQuadTree(maxSouth, maxNorth, maxWest, maxEast, minTreeMax)
	for {
		msg := <-msgChan
		msg.perf.stopAndStart(perf_tmProc)
		switch msg.op {
		case cAddOp:
			handleAdd(&msg, tree, lg)
		case cRemoveOp:
			handleRemove(&msg, tree, lg)
		case cMoveOp:
			handleMove(&msg, tree, lg)
		case cNearbyOp:
			handleNearby(&msg, tree, lg)
		}
		lg.Println(msg.perf.stopAndString())
	}
}

func handleAdd(add *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := &add.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	log.Println("User: %d \t Add Request \tmNS: %f mEW: %f", usr.id, mNS, mEW)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, addFun(usr))
	tree.Insert(mNS, mEW, usr)
}

func handleRemove(rmv *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := &rmv.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	log.Println("User: %d \t Remove Request \tmNS: %f mEW: %f", usr.id, mNS, mEW)
	deleteUsr(mNS, mEW, usr, tree)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, removeFun(usr))
}

func handleNearby(nby *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := nby.usr
	mNS, mEW := metresFromOrigin(usr.Lat, usr.Lng)
	log.Println("User: %d \t Nearby Request \t mNS %f mEW %f", usr.id, mNS, mEW)
	view := nearbyView(mNS, mEW)
	vs := []*quadtree.View{view}
	tree.Survey(vs, nearbyFun(&usr))
}

func handleMove(mv *clientMsg, tree quadtree.T, lg *log.Logger) {
	usr := &mv.usr
	nMNS, nMEW := metresFromOrigin(usr.Lat, usr.Lng)
	oMNS, oMEW := metresFromOrigin(usr.OLat, usr.OLng)
	log.Println("User: %d \t Relocate Request: \t oMNS: %f oMEW %f nMNS: %f nMEW %f", usr.id, oMNS, oMEW, nMNS, nMEW)
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
	movedView := []*quadtree.View{nView.Intersect(oView)}
	tree.Survey(movedView, movedFun(&mv.usr))
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
func nearbyFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(sNearbyOp, usr, oUsr)
		}
	}
}

// Returns a function used for alerting users that another user has been added to the system
func addFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(sAddOp, usr, oUsr)
		}
	}
}

// Returns a function used for alerting users that another user has been removed from the system
// NB: Relies on the assumption that usr is not currently present in tree
func removeFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		broadcastSend(sRemoveOp, usr, oUsr)
	}
}

// Returns a function used for alerting users that another user is going out of range and should be removed
func notVisibleFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		broadcastSend(sNotVisibleOp, usr, oUsr)
		broadcastSend(sNotVisibleOp, oUsr, usr)
	}
}

// Returns a function used for alerting users that another user has just become visible and should be added
func visibleFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(sVisibleOp, usr, oUsr)
			broadcastSend(sVisibleOp, oUsr, usr)
		}
	}
}

// Returns a function used for alerting users that another user has changed position and should be updated
func movedFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			broadcastSend(sMovedOp, usr, oUsr)
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

func broadcastSend(op serverOp, usr *user, oUsr *user) {
	perf := newPerfProfile(usr.id, usr.tId, string(op), perf_outTaskNum)
	perf.start(perf_bSend)
	msg := newServerMsg(op, usr, perf)
	oUsr.writeChan <- msg
}
