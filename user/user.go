package user

import (
	"fmt"
	"log"
	"websocket"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
)

type U struct {
	Id      string
	msgChan chan *msgdef.ServerMsg
	closeChan chan bool
}

func New() *U {
	mc := make(chan *msgdef.ServerMsg, 32)
	cc := make(chan bool, 1)
	return &U{msgChan: mc, closeChan: cc}
}

func (usr *U) WriteMsg(msg *msgdef.ServerMsg) {
	usr.msgChan<-msg
}

func (usr *U) ReceiveMsg() *msgdef.ServerMsg {
	return <-usr.msgChan
}

func (usr *U) WriteClose() {
	usr.closeChan<-true
}

func (usr *U) ReceiveClose() {
	<-usr.closeChan
	return
}

func (usr *U) ListenAndWrite(ws *websocket.Conn) {
	defer closeWS(ws, usr)
	for {
		msg := usr.ReceiveMsg()
		msgLog(usr.Id, "Server Msg", fmt.Sprintf("%v", msg))
		err := jsonutil.JSONCodec.Send(ws, msg)
		if err != nil {
			msgLog(usr.Id, "Error", err.Error())
			return
		}
		if msg.Op == msgdef.SErrorOp {
			closeWS(ws, usr)
			return
		}
	}
}

func closeWS(ws *websocket.Conn, usr *U) {
	if err := ws.Close(); err != nil {
		msgLog(usr.Id, "Error", err.Error())
	}
	usr.WriteClose()
}

func msgLog(id, title, msg string) {
	log.Printf("User: %s\t%s\t%s", id, title, msg)
}
