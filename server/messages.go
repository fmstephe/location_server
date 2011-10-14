package locserver

// ----------CLIENT MESSAGES------------
// User request operations NB: a user is synonymous with a websocket connection
type clientOp string
// Set the name of the user
const initOp = clientOp("cInit")
// Change the location of the user
const cMoveOp = clientOp("cMove")
// Query for all cNearby users
const cNearbyOp = clientOp("cNearby")

// A structure for reading messages off a websocket
type CJsonMsg struct {
	Op       clientOp
	Lat, Lng float64
	Name     string
}

// A client request to update the location of the user
type cRelocate struct {
	nMNS, nMEW float64
	oMNS, oMEW float64
	oLat, oLng float64
	usr        user
}

// A client request for all users nearby this user
type cNearby struct {
	mNS, mEW float64
	usr      user
}

// ----------SERVER MESSAGES------------
// Server reply operations
type serverOp string

// Indicates that that a user has moved out of bounds - client should remove this user
type sOutOfBounds struct {
	Op  serverOp
	Usr user
}

func newSOutOfBounds(usr *user) *sOutOfBounds {
	return &sOutOfBounds{Op: serverOp("sOutOfBounds"), Usr: *usr}
}

// Indicates that a user has been added - client should add this user
type sAdd struct {
	Op  serverOp // Always has the value "new"
	Usr user
}

func newSAdd(usr *user) *sAdd {
	return &sAdd{Op: serverOp("sAdd"), Usr: *usr}
}

// Indicates that a user has moved - client should update this user
type sMoved struct {
	Op         serverOp // Always has the value "moved"
	OLat, OLng float64
	Usr        user
}

func newSMoved(oLat, oLng float64, usr *user) *sMoved {
	return &sMoved{Op: serverOp("sMoved"), OLat: oLat, OLng: oLng, Usr: *usr}
}

//Indicates that a user has just appeared within your visible range - client should add this user
type sVisible struct {
	Op  serverOp // Always has the value "visible"
	Usr user
}

func newSVisible(usr *user) *sVisible {
	return &sVisible{Op: serverOp("sVisible"), Usr: *usr}
}

// Indicates that a user has been removed - client should remove this user
type sRemove struct {
	Op  serverOp // Always has the value "remove"
	Usr user
}

func newSRemove(usr *user) *sRemove {
	return &sRemove{Op: serverOp("sRemove"), Usr: *usr}
}

// Indicates that a user is nearby
type sNearby struct {
	Op  serverOp // Always has the value "nearby"
	Usr user
}

func newSNearby(usr *user) *sNearby {
	return &sNearby{Op: serverOp("sNearby"), Usr: *usr}
}
