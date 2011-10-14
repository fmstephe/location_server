package main

import (
	"websocket"
	"strings"
	"json"
)

func main() {
	ws, err := websocket.Dial("ws://localhost:8001/ws", "", "http://localhost:8001/")
	if err != nil {
		panic("Dial: " + err.String())
	}
	init := CJsonMsg{Op: "cInit", Lat: 1.0, Lng: 1.0, Name: "Test1"}
	cMsg := json.MarshalForHTML(init)
	println(cMsg)
	if _, err := ws.Write([]byte(cMsg)); err != nil {
		panic("Write: " + err.String())
	}
	var sMsg = make([]byte)
	if n, err := ws.Read(sMsg); err != nil {
		panic("Read: " + err.String())
	}
	println(sMsg)
}
