package main

import (
	"locserver"
	"websocket"
	"io/ioutil"
	"http"
	l4g "log4go.googlecode.com/hg"
	"flag"
	"runtime/pprof"
	"os"
	"time"
)

const index = "index.html"
const logPath = "/var/log/locserver/server.log"

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
var profiletime = flag.Int64("profiletime", 60, "The amount of time (seconds) that profiling runs for, defaults to 60")

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

func stopCPUProfiling() {
  f, _ := os.Create(*cpuprofile)
  pprof.StartCPUProfile(f)
  l4g.Info("Pofiling Started")
  time.Sleep((*profiletime)*1000000000) // One Minute of profiling
  pprof.StopCPUProfile()
  l4g.Info("Pofiling Finished")
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
	  go stopCPUProfiling()
	}
	l4g.AddFilter("file", l4g.INFO, l4g.NewFileLogWriter(logPath, false))
	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(locserver.WebsocketUser))
	go locserver.TreeManager()
	http.ListenAndServe(":8001", nil)
}
