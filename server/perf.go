package locserver

import (
	"time"
)

// ----------PERFORMANCE TRACKING------------

// Performance measured steps
//  1: User Processing
//  2: Treemanager Channel Send
//  3: Treemanager Processing
//  4: User Broadcast Channel Send
//  5: Websocket Send

type inPerf struct {
	// The user-id of the user who received the initial message for this transaction
	uId int64
	// The id of this transaction - (uId,tId) is globally unique between server restarts
	tId int64
	// Nanosecond performance timings
	// The amount of time it takes the user goroutine to unmarshal and forward a client message
	userProc int64
	// The amount of time it takes from when this message is put onto a tree manager channel and when it is taken off
	tmSend int64
	// The amount of time it takes the tree manager to process a given message
	tmProc int64
}

func newInPerf(uId, tId int64) inPerf {
	return inPerf{uId: uId, tId: tId}
}

func (p inPerf) startUserProcessing() {
	p.userProc = time.Nanoseconds()
}

func (p inPerf) endUserProcessing() {
	p.userProc = time.Nanoseconds() - p.userProc
}

func (p inPerf) startTreeManagerSend() {
	p.tmSend = time.Nanoseconds()
}

func (p inPerf) endTreeManagerSend() {
	p.tmSend = time.Nanoseconds() - p.tmSend
}

func (p inPerf) startTreeManagerProcessing() {
	p.tmProc = time.Nanoseconds()
}

func (p inPerf) endTreeManagerProcessing() {
	p.tmProc = time.Nanoseconds() - p.tmProc
}

type outPerf struct {
	// The user-id of the user who received the initial message for this transaction
	uId int64
	// The id of this transaction - (uId,tId) is globally unique between server restarts
	tId int64
	// The amount of time it takes to send a message to a user via it's writeChan channel
	bSend int64
	// The amount of time it takes for the function ws.Send(...) to complete (TODO may or may not be a useful measure - check this out)
	wSend int64
}

func newOutPerf(perf inPerf) outPerf {
	return outPerf{uId: perf.uId, tId: perf.tId}
}

func (p outPerf) startBroadcastSend() {
	p.bSend = time.Nanoseconds()
}

func (p outPerf) endBroadcastSend() {
	p.bSend = time.Nanoseconds() - p.bSend
}

func (p outPerf) startWebsocketSend() {
	p.wSend = time.Nanoseconds()
}

func (p outPerf) endWebsocketSend() {
	p.wSend = time.Nanoseconds() - p.wSend
}
