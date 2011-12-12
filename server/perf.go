package locserver

import (
	"bytes"
	"fmt"
	"time"
)

// Incoming performance tasks
const perf_inTaskNum = 3
const perf_userProc = "userProc"
const perf_tmSend = "tmSend"
const perf_tmProc = "tmProc"

// Outgoing performance tasks
const perf_outTaskNum = 3
const perf_bSend = "bSend"
const perf_bProc = "bProc"
const perf_wSend = "wSend"

// ----------PERFORMANCE TRACKING------------

// Performance measured steps
//  1: User Processing
//  2: Treemanager Channel Send
//  3: Treemanager Processing
//  4: User Broadcast Channel Send
//  5: Websocket Send

type perfProfiler interface {
	perfProfile() *perfProfile
}

type perfProfile struct {
	// A, preferably unique, name for this performance profile
	pName string
	// Nanosecond task performance timings
	timings []perfUnit
}

// A performance profile unit - represents the timing of a specific task
type perfUnit struct {
	taskName string
	time     int64
}

func newPerfProfile(uId, tId int64, op string, taskNum int) *perfProfile {
	t := make([]perfUnit, 0, taskNum)
	return &perfProfile{pName: fmt.Sprintf("%d:%d:%s", uId, tId, op), timings: t}
}

func (p *perfProfile) start(taskName string) {
	u := perfUnit{taskName: taskName, time: time.Now()}
	p.timings = append(p.timings, u)
}

func (p *perfProfile) stop() {
	last := &p.timings[len(p.timings)-1]
	last.time = time.Now().Sub(last.time)
}

func (p *perfProfile) stopAndStart(taskName string) {
	p.stop()
	p.start(taskName)
}

func (p *perfProfile) stopAndString() string {
	p.stop()
	buf := bytes.NewBufferString("perf-" + p.pName + "\t")
	for i := range p.timings {
		unit := &p.timings[i]
		fmt.Fprintf(buf, "%s %10.6f\t", unit.taskName, toMilli(unit.time)) // Work out what to do with the error
	}
	return buf.String()
}

func toMilli(nano int64) float64 {
	short := int32(nano)
	return float64(short) / 1000000
}
