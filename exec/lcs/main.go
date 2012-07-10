package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"github.com/fmstephe/location_server/locserver"
	"github.com/fmstephe/location_server/logutil"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

const logPath = "/var/log/locserver/server.log"

var iFile []byte
var minTreeMax *int64 = flag.Int64("treeSize", 1000, "The initialisation size of the quadtree")
var trackMovement *bool = flag.Bool("m", false, "Broadcast fine grained movement of users")
var threads *int = flag.Int("t", 1, "The number of threads available to the runtime")

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*threads)
}

func main() {
	logutil.ServerStarted("Location")
	http.Handle("/loc", websocket.Handler(locserver.HandleLocationService))
	locserver.StartTreeManager(*minTreeMax, *trackMovement)
	http.ListenAndServe(":8002", nil)
}
