package jsonutil

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"github.com/fmstephe/location_server/logutil"
	"html"
)

// Unmarshals a websocket message into msgi as JSON.
func UnmarshalAndLog(tId uint, uId string, ws *websocket.Conn, msg interface{}) error {
	var data string
	if err := websocket.Message.Receive(ws, &data); err != nil {
		return err
	}
	logutil.Log(tId, uId, data)
	if err := json.Unmarshal([]byte(data), msg); err != nil {
		return err
	}
	return nil
}

func UnmarshalAndProcess(tId uint, uId string, ws *websocket.Conn, msg interface{}, processFunc func() error) error {
	UnmarshalAndLog(tId, uId, ws, msg)
	return processFunc()
}

// 
func SanitiseJSON(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return html.EscapeString(s)
	}
	sanitiseJSON(v)
	return v
}

func sanitiseJSON(parent interface{}) {
	switch parent.(type) {
	case map[string]interface{}:
		parentMap := parent.(map[string]interface{})
		for k, child := range parentMap {
			switch child.(type) {
			case string:
				parentMap[k] = html.EscapeString(child.(string))
			case map[string]interface{}:
				SanitiseJSON(child)
			}
		}
	}
}
