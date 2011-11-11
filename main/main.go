package main

import (
	"location_server/server"
	"websocket"
	"io/ioutil"
	"net/http"
	"log"
	"os"
	_ "net/http/pprof"
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

// TODO This is poorly - improve
func initLog() *log.Logger {
	logFile, err := os.OpenFile(logPath, os.O_WRONLY , 0666)
	if err != nil {
		os.Create(logPath)
		logFile, _ = os.OpenFile(logPath, os.O_WRONLY , 0666)
	}
	return log.New(logFile, "", log.Lmicroseconds)
  }


func main() {
	lg := initLog()
	lg.Println("Location Server Started")
	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager(minTreeMax, lg)
	http.ListenAndServe(":8001", nil)
}
