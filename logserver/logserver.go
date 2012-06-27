package logserver

import (
	"code.google.com/p/go.net/websocket"
	"location_server/msgutil/msgwriter"
	"location_server/msgutil/jsonutil"
)

func HandleLogService(ws *websocket.Conn) {
	var tId uint
	uId := "N/A"
	msgWriter := msgwriter.New(ws)
	for {
		var msg interface{}
		if err := jsonutil.UnmarshalAndLog(tId, uId, ws, msg); err != nil {
			msgWriter.ErrorAndClose(tId, uId, err.Error())
			return
		}
	}
}
