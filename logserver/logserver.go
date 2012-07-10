package logserver

import (
	"code.google.com/p/go.net/websocket"
	"github.com/fmstephe/location_server/msgutil/jsonutil"
	"github.com/fmstephe/location_server/msgutil/msgwriter"
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
