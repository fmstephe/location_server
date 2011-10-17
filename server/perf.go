package locserver

import (
	"time"
	l4g "log4go.googlecode.com/hg"
)

// ----------PERFORMANCE TRACKING------------

// Performance measured steps
//  1: User Processing
//  2: Treemanager Channel Send
//  3: Treemanager Processing
//  4: User Broadcast Channel Send
//  5: Websocket Send

type inPerf struct {
	// The operation for this transaction
	op clientOp
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

func newInPerf(op clientOp, uId, tId int64) *inPerf {
	return &inPerf{op: op, uId: uId, tId: tId}
}

func (p *inPerf) beginUserProc() {
	p.userProc = time.Nanoseconds()
}

func (p *inPerf) beginTmSend() {
	p.userProc = time.Nanoseconds() - p.userProc
	p.tmSend = time.Nanoseconds()
}

func (p *inPerf) beginTmProc() {
	p.tmSend = time.Nanoseconds() - p.tmSend
	p.tmProc = time.Nanoseconds()
}

func (p *inPerf) finishAndLog() {
	p.tmProc = time.Nanoseconds() - p.tmProc
	fStr := "Transaction: %d:%d \tClient Op %s\tUser Processing %d\tTree Manager Msg Send %d\tTree Manager Processing %d"
	l4g.Info(fStr, p.uId, p.tId, p.op, p.userProc, p.tmSend, p.tmProc)
}

type outPerfer interface {
	getOutPerf() *outPerf
}

type outPerf struct {
	// The operation for this outbound message
	op serverOp
	// The user-id of the user who received the initial message for this transaction
	uId int64
	// The id of this transaction - (uId,tId) is globally unique between server restarts
	tId int64
	// The amount of time it takes to send a message to a user via it's writeChan channel
	bSend int64
	// The amount of time it takes for the function ws.Send(...) to complete (TODO may or may not be a useful measure - check this out)
	wSend int64
}

func newOutPerf(op serverOp, perf inPerf) outPerf {
	return outPerf{uId: perf.uId, tId: perf.tId}
}

func (p *outPerf) beginBSend() {
	p.bSend = time.Nanoseconds()
}

func (p *outPerf) beginWSend() {
	p.wSend = time.Nanoseconds()
}

func (p *outPerf) finishAndLog() {
	p.wSend = time.Nanoseconds() - p.wSend
	fStr := "Transaction: %d:%d \tServer Op %s\tUser Processing %d\tBroadcast Send %d\tWebsocket Send %d"
	l4g.Info(fStr, p.uId, p.tId, p.op, p.bSend, p.wSend)
}
