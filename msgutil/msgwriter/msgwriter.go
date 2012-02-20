package msgwriter

import (
	"websocket"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
)

type W struct {
	ws *websocket.Conn
	msgChan chan *msgdef.ServerMsg
	closeChan chan bool
}

func New(ws *websocket.Conn) *W {
	msgChan := make(chan *msgdef.ServerMsg, 32)
	closeChan := make(chan bool, 1)
	msgWriter := &W{ws: ws, msgChan: msgChan, closeChan: closeChan}
	go msgWriter.listenAndWrite()
	return msgWriter
}

func (msgWriter *W) WriteMsg(msg *msgdef.ServerMsg) {
	msgWriter.msgChan<-msg
}

func (msgWriter *W) ErrorAndClose(tId uint, errMsg string) {
	//msgLog(usr.Id, "Connection Terminated", errMsg)
	msgWriter.WriteMsg(msgdef.NewServerError(errMsg))
	<-msgWriter.closeChan
	//msgLog(usr.Id, "Close Confirmation Received", "")
}

func (msgWriter *W) listenAndWrite() {
	defer msgWriter.closeWS()
	for {
		msg := <-msgWriter.msgChan
		// TODO log the message here
		err := jsonutil.JSONCodec.Send(msgWriter.ws, msg)
		if err != nil {
			// TODO log the error here
			return
		}
		if msg.Op == msgdef.SErrorOp {
			return
		}
	}
}

func (msgWriter *W) closeWS() {
	if err := msgWriter.ws.Close(); err != nil {
		//msgLog("Error closing websocket connection: ", err.Error())
	}
	msgWriter.closeChan<-true
}
