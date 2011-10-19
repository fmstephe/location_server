package main

import (
	"flag"
	"time"
	"locserver"
	"websocket"
	"json"
	"rand"
)

const one_second = 1000000000

var workers = flag.Int("w", 1, "The number of workers to be spawned")

func main() {
	flag.Parse()
	params := wanderParams
	sleepTime := int64(one_second)
	for i := 0; i < *workers; i++ {
		lat, lng, init, nxtPos := params()
		go run_test(lat, lng, init, nxtPos, sleepTime)
	}
	lat, lng, init, nxtPos := params()
	run_test(lat, lng, init, nxtPos, sleepTime)
}

func jmpParams() (lat, lng float64, init *locserver.CJsonMsg, nxtPos func(float64, float64) (float64, float64)) {
	lat = 1.0
	lng = 1.0
	init = locserver.TestMsg("cInit", lat, lng, "name")
	nxtPos = jmpPos
	return
}

func jmpPos(lat, lng float64) (nLat, nLng float64) {
	nLat = float64(lat + 1)
	nLng = float64(lng + 1)
	// A kind of mod
	if nLat > 90 {
		nLat = -90
	}
	if nLng > 180 {
		nLng = -180
	}
	return
}

func wanderParams() (lat, lng float64, init *locserver.CJsonMsg, nxtPos func(float64, float64) (float64, float64)) {
	lat = rand.Float64() + 0.5
	lng = rand.Float64() + 0.5
	init = locserver.TestMsg("cInit", lat, lng, "wanderer")
	nxtPos = wanderPos
	return
}

func wanderPos(lat, lng float64) (nLat, nLng float64) {
	nLat = lat + (rand.Float64() * 0.02) - 0.01
	nLng = lng + (rand.Float64() * 0.02) - 0.01
	// This is awful
	if nLat > 1.5 {
		nLat = 1.5
	}
	if nLat < 0.5 {
		nLat = 0.5
	}
	if nLng > 1.5 {
		nLng = 1.5
	}
	if nLng < 0.5 {
		nLng = 0.5
	}
	return
}

func run_test(lat, lng float64, init *locserver.CJsonMsg, nxtPos func(float64, float64) (float64, float64), sleepTime int64) {
	ws := doDial()
	go eatMsgs(ws)
	marshalAndSend(init, ws)
	i := 0
	for {
		time.Sleep(sleepTime)
		i++
		lat, lng = nxtPos(lat, lng)
		cMsg := locserver.TestMsg("cMove", lat, lng, "worthless")
		cMsgA, err := json.MarshalForHTML(cMsg)
		if err != nil {
			return
		}
		if _, err := ws.Write([]byte(cMsgA)); err != nil {
			return
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

func doDial() *websocket.Conn {
	ws, err := websocket.Dial("ws://localhost:8001/ws", "", "http://localhost:8001/")
	if err != nil {
		panic("Dial: " + err.String())
	}
	return ws
}

func marshalAndSend(msg interface{}, ws *websocket.Conn) {
	msgA, err := json.MarshalForHTML(msg)
	if err != nil {
		panic("Unmarshall: " + err.String())
	}
	if _, err := ws.Write([]byte(msgA)); err != nil {
		panic("Write: " + err.String())
	}
}
