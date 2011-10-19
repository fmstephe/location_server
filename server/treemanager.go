package locserver

import (
	"quadtree"
	l4g "log4go.googlecode.com/hg"
)

const (
	nearbyMetresNS = 1000.0
	nearbyMetresEW = 1000.0
)

//
// Single Threaded Tree Manager Code 
//
var msgChan = make(chan interface{}, 255) // Global Channel for all requests

func TreeManager() {
	tree := quadtree.NewQuadTree(maxSouth, maxNorth, maxWest, maxEast)
	for {
		msg := <-msgChan
		switch t := msg.(type) {
			case cAdd: 
			add := msg.(cAdd)
			handleAdd(&add, tree)
			case cRemove: 
			rmv := msg.(cRemove)
			handleRemove(&rmv, tree)
			case cMove : 
			mv := msg.(cMove)
			handleMove(&mv, tree)
			case cNearby : 
			nby := msg.(cNearby)
			handleNearby(&nby, tree)
		}
	}
}

func handleAdd(add *cAdd, tree quadtree.QuadTree) {
	add.perf.beginTmProc()
	usr := &add.usr
	l4g.Info("User: %d \t Add Request \tmNS: %f mEW: %f", usr.id, usr.mNS, usr.mEW)
	mNS := usr.mNS
	mEW := usr.mEW
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, addFun(usr))
	tree.Insert(mNS, mEW, usr)
	add.perf.finishAndLog()
}

func handleRemove(rmv *cRemove, tree quadtree.QuadTree) {
	rmv.perf.beginTmProc()
	usr := &rmv.usr
	l4g.Info("User: %d \t Remove Request", usr.id)
	mNS := usr.mNS
	mEW := usr.mEW
	deleteUsr(mNS, mEW, usr, tree)
	vs := []*quadtree.View{nearbyView(mNS, mEW)}
	tree.Survey(vs, removeFun(usr))
	rmv.perf.finishAndLog()
}

func handleNearby(nby *cNearby, tree quadtree.QuadTree) {
	nby.perf.beginTmProc()
	l4g.Info("User: %d \t Nearby Request \t mNS %f mEW %f", nby.usr.id, nby.mNS, nby.mEW)
	usr := nby.usr
	view := nearbyView(usr.mNS, usr.mEW)
	vs := []*quadtree.View{view}
	tree.Survey(vs, nearbyFun(&usr))
	nby.perf.finishAndLog()
}

func handleMove(mv *cMove, tree quadtree.QuadTree) {
	mv.perf.beginTmProc()
	l4g.Info("User: %d \t Relocate Request: \t oMNS: %f oMEW %f nMNS: %f nMEW %f", mv.usr.id, mv.oMNS, mv.oMEW, mv.nMNS, mv.nMEW)
	usr := &mv.usr
	deleteUsr(mv.oMNS, mv.oMEW, usr, tree)
	tree.Insert(mv.nMNS, mv.nMEW, usr)
	nView := nearbyView(mv.nMNS, mv.nMEW)
	oView := nearbyView(mv.oMNS, mv.oMEW)
	// Alert out of bounds users
	nvViews := oView.Subtract(nView)
	tree.Survey(nvViews, notVisibleFun(mv))
	// Alert newly visible users
	vViews := nView.Subtract(oView)
	tree.Survey(vViews, visibleFun(mv))
	// Alert watching users of the relocation
	// movedView := []*quadtree.View{nView.Intersect(oView)}
	// tree.Survey(movedView, movedFun(mv))
	mv.perf.finishAndLog()
}

// Deletes usr from tree at the given coords
func deleteUsr(mNS, mEW float64, usr *user, tree quadtree.QuadTree) {
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
			sNby := newSNearby(oUsr)
			sNby.perf.beginBSend()
			usr.writeChan <- sNby
		}
	}
}

// Returns a function used for alerting users that another user has been added to the system
func addFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !usr.eq(oUsr) {
			sAdd := newSAdd(usr)
			sAdd.perf.beginBSend()
			oUsr.writeChan <- sAdd
		}
	}
}

// Returns a function used for alerting users that another user has been removed from the system
// NB: Relies on the assumption that usr is not currently present in tree
func removeFun(usr *user) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		sRmv := newSRemove(usr)
		sRmv.perf.beginBSend()
		oUsr.writeChan <- sRmv
	}
}

// Returns a function used for alerting users that another user is going out of range and should be removed
func notVisibleFun(mv *cMove) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		// Send nv to other user
		uOob := newSNotVisible(&mv.usr)
		uOob.perf.beginBSend()
		oUsr.writeChan <- uOob
		// Send nv to moving user
		ouOob := newSNotVisible(&mv.usr)
		ouOob.perf.beginBSend()
		mv.usr.writeChan <- ouOob
	}
}

// Returns a function used for alerting users that another user has just become visible and should be added
func visibleFun(mv *cMove) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !mv.usr.eq(oUsr) {
			uVsb := newSVisible(&mv.usr)
			uVsb.perf.beginBSend()
			oUsr.writeChan <- uVsb
			ouVsb := newSVisible(oUsr)
			ouVsb.perf.beginBSend()
			mv.usr.writeChan <- ouVsb
		}
	}
}

// Returns a function used for alerting users that another user has changed position and should be updated
func movedFun(mv *cMove) func(mNS, mEW float64, e interface{}) {
	return func(mNS, mEW float64, e interface{}) {
		oUsr := e.(*user)
		if !mv.usr.eq(oUsr) {
			mvd := newSMoved(mv.oLat, mv.oLng, &mv.usr)
			mvd.perf.beginBSend()
			oUsr.writeChan <- mvd
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
