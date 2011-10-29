package main

import (
	"locserver"
	"websocket"
	"io/ioutil"
	"http"
	_ "http/pprof"
)

const index = "index.html"
const logPath = "/var/log/locserver/server.log"

var minTreeMax = int64(1000000)

//
// Static index HTML page serving function
//
func indexHandler(w http.ResponseWriter, r *http.Request) {
	iFile, err := ioutil.ReadFile(index)
	if err != nil {
		return
	}
	w.Write(iFile)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager(minTreeMax)
	http.ListenAndServe(":8001", nil)
}
