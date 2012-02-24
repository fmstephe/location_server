package msgwriter

import (
	"fmt"
	"location_server/logutil"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
	"websocket"
)

type W struct {
	ws        *websocket.Conn
	msgChan   chan *msgdef.ServerMsg
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
	msgWriter.msgChan <- msg
}

func (msgWriter *W) ErrorAndClose(tId uint, uId, errMsg string) {
	logutil.Log(tId, uId, "Connection Terminated: "+errMsg)
	msgWriter.WriteMsg(msgdef.NewServerError(errMsg))
	<-msgWriter.closeChan
	logutil.Log(tId, uId, "Close Confirmation Received")
}

func (msgWriter *W) listenAndWrite() {
	defer msgWriter.sendClose()
	for {
		msg := <-msgWriter.msgChan
		logutil.Log(msg.TransactionId(), msg.UserId(), fmt.Sprintf("%v", msg))
		err := jsonutil.JSONCodec.Send(msgWriter.ws, msg)
		if err != nil {
			logutil.Log(msg.TransactionId(), msg.UserId(), err.Error())
			return
		}
		if msg.Op == msgdef.SErrorOp {
			return
		}
	}
}

func (msgWriter *W) sendClose() {
	msgWriter.closeChan <- true
}