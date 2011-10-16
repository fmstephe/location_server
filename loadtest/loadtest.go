package main

import (
	"strconv"
	"time"
	"locserver"
	"websocket"
	"json"
)

const one_second = 1000000000

func main() {
	for i := 0; i < 1000; i++ {
		go run_test("Test_" + strconv.Itoa(i))
		time.Sleep(one_second / 8)
	}
	run_test("Test_End")
}

func run_test(name string) {
	ws, err := websocket.Dial("ws://localhost:8001/ws", "", "http://localhost:8001/")
	go eatMsgs(ws)
	if err != nil {
		panic("Dial: " + err.String())
	}
	init := locserver.CJsonMsg{Op: "cInit", Lat: 1.0, Lng: 1.0, Name: name}
	initA, err := json.MarshalForHTML(init)
	//println(string(initA))
	if _, err := ws.Write([]byte(initA)); err != nil {
		panic("Write: " + err.String())
	}
	i := 0
	for {
		time.Sleep(one_second / 16)
		i++
		lat := float64(i % 90)
		lng := float64(i % 180)
		cMsg := locserver.CJsonMsg{Op: "cMove", Lat: lat, Lng: lng}
		cMsgA, err := json.MarshalForHTML(cMsg)
		if err != nil {
			panic("Write: " + err.String())
		}
		//println(string(cMsgA))
		if _, err := ws.Write([]byte(cMsgA)); err != nil {
			panic("Write: " + err.String())
		}
	}
}

func eatMsgs(ws *websocket.Conn) {
	var sMsg = make([]byte, 256, 256)
	for {
		n, err := ws.Read(sMsg)
		if err != nil {
			panic("Read: " + err.String())
		}
		println(string(sMsg[:n]))
	}
}
