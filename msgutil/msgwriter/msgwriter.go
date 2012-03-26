package msgwriter

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"location_server/logutil"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
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
	msgWriter.WriteMsg(msgdef.NewServerError(tId, uId, errMsg))
	<-msgWriter.closeChan
	logutil.Log(tId, uId, "Close Confirmation Received")
}

func (msgWriter *W) listenAndWrite() {
	defer msgWriter.sendClose()
	for {
		sMsg := <-msgWriter.msgChan
		msg := sMsg.Msg
		logutil.Log(sMsg.TId, sMsg.UId, fmt.Sprintf("Server Sent: %v", msg))
		err := jsonutil.JSONCodec.Send(msgWriter.ws, msg)
		if err != nil {
			logutil.Log(sMsg.TId, sMsg.UId, err.Error())
			return
		}
		if _, ok := msg.(*msgdef.SErrorMsg); ok {
			logutil.Log(sMsg.TId, sMsg.UId, "Error message received - Shutting Down")
			return
		}
	}
}

func (msgWriter *W) sendClose() {
	msgWriter.closeChan <- true
}
