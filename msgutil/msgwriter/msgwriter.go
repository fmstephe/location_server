package msgwriter

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/fmstephe/location_server/logutil"
	"github.com/fmstephe/location_server/msgutil/msgdef"
)

// This message instructs the message writer to shutdown after writing errMsg back to its websocket
type shutdown struct {
	closeChan chan bool
	errMsg    *msgdef.ServerMsg
}

// A message writer listens on a channel for messages to write to a websocket
// A message writer listen until it receives an error message
// then it will write the error message to the websocket  and terminate
type W struct {
	ws           *websocket.Conn
	msgChan      chan *msgdef.ServerMsg
	shutdownChan chan *shutdown
}

// Creates and returns a new message writer
// Starts a goroutine listening for incoming messages
func New(ws *websocket.Conn) *W {
	msgChan := make(chan *msgdef.ServerMsg, 32)
	shutdownChan := make(chan *shutdown, 1)
	msgWriter := &W{ws: ws, msgChan: msgChan, shutdownChan: shutdownChan}
	go msgWriter.listenAndWriteback()
	return msgWriter
}

// Asks the message writer to write msg back to its websocket
func (msgWriter *W) WriteMsg(msg *msgdef.ServerMsg) {
	msgWriter.msgChan <- msg
}

// Asks the message writer to write the error message to its websocket  and terminate
// This function waits on a message from the closeChan to ensure that
func (msgWriter *W) ErrorAndClose(tId uint, uId, errMsg string) {
	logutil.Log(tId, uId, "Connection Terminated: "+errMsg)
	closeChan := make(chan bool, 1)
	sd := &shutdown{closeChan, msgdef.NewServerError(tId, uId, errMsg)}
	msgWriter.shutdownChan <- sd
	<-closeChan
	logutil.Log(tId, uId, "Close Confirmation Received")
}

// Loops listening for messages to write back to msgWriter's websocket
// There are two possible messages to receive
// 1: Server message
//	The contents of the server message is marshalled and written back to the websocket
//	If an error occurs the error is logged and written back to the websocket
//	Loop continues
// 2: shutdown message
//	The contents of the shutdown message's server message is marshalled and written back to the websocket
//	If an error occurs the error is logged and written back to the websocket
//	shutdown messages closeChan is sent a value, allowing the sender to unblock
//	Loop terminates
func (msgWriter *W) listenAndWriteback() {
	for {
		var sMsg *msgdef.ServerMsg
		var closeChan chan bool
		select {
		case sMsg = <-msgWriter.msgChan:
		case sd := <-msgWriter.shutdownChan:
			sMsg = sd.errMsg
			closeChan = sd.closeChan
		}
		if err := writeback(msgWriter.ws, sMsg); err != nil {
			logutil.Log(sMsg.TId, sMsg.UId, err.Error())
		}
		if closeChan != nil {
			logutil.Log(sMsg.TId, sMsg.UId, "Error message received - Shutting Down")
			closeChan <- true
			return
		}
	}
}

// Writes the contents of sMsg back to ws, returning any errors encountered
func writeback(ws *websocket.Conn, sMsg *msgdef.ServerMsg) error {
	msg := sMsg.Msg
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	logutil.Log(sMsg.TId, sMsg.UId, fmt.Sprintf("Server Sent: %s", bytes))
	_, err = ws.Write(bytes)
	return err
}
