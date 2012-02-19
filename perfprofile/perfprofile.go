package perfprofile

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
	Profile() *P
}

type P struct {
	pName string // A, preferably unique, name for this performance profile
	timings []perfUnit // Nanosecond task performance timings
}

// A performance profile unit - represents the timing of a specific task
type perfUnit struct {
	taskName string
	time     int64
}

func New(name string, taskNum int) *P {
	t := make([]perfUnit, 0, taskNum)
	return &P{pName: name, timings: t}
}

// Utility method to assist common profile naming pattern
func ProfileName(tId uint, uId, msg string) string {
	return fmt.Sprintf("%d:%s:%s", tId, uId, msg)
}

func (p *P) Start(taskName string) {
	u := perfUnit{taskName: taskName, time: time.Now().UnixNano()}
	p.timings = append(p.timings, u)
}

func (p *P) Stop() {
	last := &p.timings[len(p.timings)-1]
	last.time = time.Now().UnixNano() - last.time
}

func (p *P) StopAndStart(taskName string) {
	p.Stop()
	p.Start(taskName)
}

func (p *P) StopAndString() string {
	p.Stop()
	return p.String()
}

func (p *P) String() string {
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
