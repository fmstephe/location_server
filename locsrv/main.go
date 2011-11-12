package main

import (
	"flag"
	"io/ioutil"
	"location_server/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"websocket"
)

const index = "index.html"
const logPath = "/var/log/locserver/server.log"

var minTreeMax = int64(1000000)
var iFile []byte
var trackMovement *bool = flag.Bool("m", false, "Broadcast fine grained movement of users")

func init() {
	println("index request")
	var err error
	iFile, err = ioutil.ReadFile(index)
	if err != nil {
		panic("Unable to initialise index.html")
	}
}

//
// Static index HTML page serving function
//
func indexHandler(w http.ResponseWriter, r *http.Request) {
	println("index request")
	w.Write(iFile)
}

// TODO This is poorly - improve
func initLog() *log.Logger {
	logFile, err := os.OpenFile(logPath, os.O_WRONLY, 0666)
	if err != nil {
		os.Create(logPath)
		logFile, _ = os.OpenFile(logPath, os.O_WRONLY, 0666)
	}
	return log.New(logFile, "", log.Lmicroseconds)
}

func main() {
	flag.Parse()
	lg := initLog()
	lg.Println("Location Server Started")
	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager(minTreeMax, *trackMovement, lg)
	http.ListenAndServe(":8001", nil)
}
