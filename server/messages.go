package locserver

// ----------CLIENT MESSAGES------------
// User request operations NB: a user is synonymous with a websocket connection
type clientOp string
// Set the name of the useri and store user location in the tree
const initOp = clientOp("cInit")
// Change the location of the user
const cMoveOp = clientOp("cMove")
// Query for all cNearby users
const cNearbyOp = clientOp("cNearby")
// Remove user from the tree - result of closed connection
const cRemoveOp = clientOp("cRemove")

// A structure for reading messages off a websocket
type CJsonMsg struct {
	Op       clientOp
	Lat, Lng float64
	Name     string
	perf     *inPerf
}

// Usually CJsonMsg structs are built by unmarshal - but this is hard for external test packages
func TestMsg(op clientOp, lat, lng float64, name string) *CJsonMsg {
	return &CJsonMsg{Op: op, Lat: lat, Lng: lng, Name: name}
}

type cAdd struct {
	usr  user
	perf inPerf
}

func newCAdd(usr *user, perf *inPerf) cAdd {
	lPerf := *perf
	lPerf.op = initOp
	return cAdd{usr: *usr, perf: lPerf}
}

type cRemove struct {
	usr  user
	perf inPerf
}

func newCRemove(usr *user, perf *inPerf) cRemove {
	lPerf := *perf
	lPerf.op = cRemoveOp
	return cRemove{usr: *usr, perf: lPerf}
}

// A client request to update the location of the user
type cMove struct {
	nMNS, nMEW float64
	oMNS, oMEW float64
	oLat, oLng float64
	usr        user
	perf       inPerf
}

func newCMove(oLat, oLng, nLat, nLng float64, usr *user, perf *inPerf) cMove {
	mv := new(cMove)
	oMNS, oMEW := metresFromOrigin(oLat, oLng)
	nMNS, nMEW := metresFromOrigin(nLat, nLng)
	// The new metre coords
	mv.nMNS = nMNS
	mv.nMEW = nMEW
	// The old metre coords
	mv.oMNS = oMNS
	mv.oMEW = oMEW
	// The old degree coords
	mv.oLat = oLat
	mv.oLng = oLng
	// Add usr
	mv.usr = *usr
	mv.perf = *perf
	return *mv
}

// A client request for all users nearby this user
type cNearby struct {
	mNS, mEW float64
	usr      user
	perf     inPerf
}

func newCNearby(lat, lng float64, usr *user, perf *inPerf) cNearby {
	nby := new(cNearby)
	nby.mNS, nby.mEW = metresFromOrigin(lat, lng)
	nby.usr = *usr
	lPerf := *perf
	lPerf.op = cNearbyOp
	nby.perf = lPerf
	return *nby
}

// ----------SERVER MESSAGES------------
// Server reply operations
type serverOp string

// A server message which contains only a serverOp and a user
type sUserMsg struct {
	Op   serverOp
	Usr  user
	perf outPerf
}

func (msg *sUserMsg) getOutPerf() *outPerf {
	return &msg.perf
}

func (msg *sUserMsg) initPerf() {
	msg.perf.op = msg.Op
}

// Indicates that a user has been added - client should add this user
func newSAdd(usr *user) *sUserMsg {
	m := &sUserMsg{Op: serverOp("sAdd"), Usr: *usr}
	m.initPerf()
	return m
}

//Indicates that a user has just appeared within your visible range - client should add this user
func newSVisible(usr *user) *sUserMsg {
	v := &sUserMsg{Op: serverOp("sVisible"), Usr: *usr}
	v.initPerf()
	return v
}

// Indicates that that a user has moved out of bounds - client should remove this user
func newSNotVisible(usr *user) *sUserMsg {
	n := &sUserMsg{Op: serverOp("sNotVisible"), Usr: *usr}
	n.initPerf()
	return n
}

// Indicates that a user has been removed - client should remove this user
func newSRemove(usr *user) *sUserMsg {
	r := &sUserMsg{Op: serverOp("sRemove"), Usr: *usr}
	r.initPerf()
	return r
}

// Indicates that a user is nearby
func newSNearby(usr *user) *sUserMsg {
	n := &sUserMsg{Op: serverOp("sNearby"), Usr: *usr}
	n.initPerf()
	return n
}

// Indicates that a user has moved - client should update this user
type sMoved struct {
	Op         serverOp // Always has the value "moved"
	OLat, OLng float64
	Usr        user
	perf       outPerf
}

func newSMoved(oLat, oLng float64, usr *user) *sMoved {
	op := serverOp("sMoved")
	m := &sMoved{Op: op, OLat: oLat, OLng: oLng, Usr: *usr}
	m.perf.op = op
	return m
}

func (mvd *sMoved) getOutPerf() *outPerf {
	return &mvd.perf
}
