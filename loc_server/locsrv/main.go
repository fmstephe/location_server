package main

import (
	"flag"
	"io/ioutil"
	"location_server/locserver/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"websocket"
	"runtime"
)

const index = "index.html"
const logPath = "/var/log/locserver/server.log"

var iFile []byte
var minTreeMax *int64 = flag.Int64("treeSize", 1000, "The initialisation size of the quadtree")
var trackMovement *bool = flag.Bool("m", false, "Broadcast fine grained movement of users")
var threads *int = flag.Int("t", 1, "The number of threads available to the runtime")

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*threads)
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
		logFile, err = os.OpenFile(logPath, os.O_WRONLY, 0666)
		if err != nil {
			panic(err.Error())
		}
	}
	return log.New(logFile, "", log.Lmicroseconds)
}

func main() {
	lg := initLog()
	lg.Println("Location Server Started")
	http.HandleFunc("/", indexHandler)
	http.Handle("/loc", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager(*minTreeMax, *trackMovement, lg)
	http.ListenAndServe(":8002", nil)
}
