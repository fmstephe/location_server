package user

import (
	"location_server/msgdef"
)

// A user
type U struct {
	Id         string          // Unique identifier for this user
	tId        int64           // A changing identifier used to tag each message with a transaction id
	OLat, OLng float64         // Previous position of this user
	Lat, Lng   float64         // Current position of this user
	writeChan chan *msgdef.ServerMsg
	closeChan chan bool
}

func New() *U {
	wc := make(chan *msgdef.ServerMsg, 32)
	cc := make(chan bool, 1)
	return &U{writeChan: wc, closeChan: cc}
}

func (usr *U) WriteMsg(msg *msgdef.ServerMsg) {
	usr.writeChan<-msg
}

func (usr *U) ReceiveMsg() *msgdef.ServerMsg {
	return <-usr.writeChan
}

func (usr *U) WriteClose() {
	usr.closeChan<-true
}

func (usr *U) ReceiveClose() {
	<-usr.closeChan
	return
}

func (usr *U) TransactionId () int64 {
	return usr.tId
}

func (usr *U) NewTransactionId() int64 {
	usr.tId++
	return usr.tId
}

func (usr *U) Eq(oUsr *U) bool {
	if usr == nil || oUsr == nil {
		return false
	}
	return usr.Id == oUsr.Id
}
