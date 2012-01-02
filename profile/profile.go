package profile

import (
	"bytes"
	"fmt"
	"time"
)

// ----------PERFORMANCE TRACKING------------

// Performance measured steps
//  1: User Processing
//  2: Treemanager Channel Send
//  3: Treemanager Processing
//  4: User Broadcast Channel Send
//  5: Websocket Send

type Profiler interface {
	Profile() P
}

type P interface {
	Start(string)
	Stop()
	StopAndStart(string)
	StopAndString()string
	String()string
}

type profile struct {
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

func New(uId string, tId int64, op string, taskNum int) *profile {
	t := make([]perfUnit, 0, taskNum)
	return &profile{pName: fmt.Sprintf("%s:%d:%s", uId, tId, op), timings: t}
}

func (p *profile) Start(taskName string) {
	u := perfUnit{taskName: taskName, time: time.Now().UnixNano()}
	p.timings = append(p.timings, u)
}

func (p *profile) Stop() {
	last := &p.timings[len(p.timings)-1]
	last.time = time.Now().UnixNano() - last.time
}

func (p *profile) StopAndStart(taskName string) {
	p.Stop()
	p.Start(taskName)
}

func (p *profile) StopAndString() string {
	p.Stop()
	return p.String()
}

func (p *profile) String() string {
	buf := bytes.NewBufferString("perf-" + p.pName + "\t")
	for i := range p.timings {
		unit := &p.timings[i]
		fmt.Fprintf(buf, "%s %10.6f\t", unit.taskName, toMilli(unit.time))
	}
	return buf.String()
}

func toMilli(nano int64) float64 {
	short := int32(nano)
	return float64(short) / 1000000
}
