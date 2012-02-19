package jsonutil

import (
	"websocket"
	"encoding/json"
)

func jsonMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = json.MarshalForHTML(v)
	return msg, websocket.TextFrame, err
}

func jsonUnmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return json.Unmarshal(msg, v)
}

var JSONCodec = websocket.Codec{jsonMarshal,jsonUnmarshal}
