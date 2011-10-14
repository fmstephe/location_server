package main

import (
	"locserver"
	"websocket"
	"io/ioutil"
	"http"
	l4g "log4go.googlecode.com/hg"
)

const index = "index.html"
const logPath = "/var/log/locserver/server.log"

//
// Static index HTML page serving function
//
func indexHandler(w http.ResponseWriter, r *http.Request) {
	l4g.Info("Default HTML demo retrieved")
	iFile, err := ioutil.ReadFile(index)
	if err != nil {
		return
	}
	w.Write(iFile)
}

func main() {
	l4g.AddFilter("file", l4g.INFO, l4g.NewFileLogWriter(logPath, false))
	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager()
	http.ListenAndServe(":8001", nil)
}
