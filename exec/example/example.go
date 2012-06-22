package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fmstephe/simpleid"
	"location_server/locserver"
	"location_server/logutil"
	"location_server/msgserver"
	"location_server/msgutil/msgdef"
	"net/http"
	"os"
)

var port = flag.Int("port", 80, "Sets the port the server will attach to")

var idMaker = simpleid.NewIdMaker()

// Provides a unique id in a nice json msg (Op: "sIdOp", Id: $$$)
func idProvider(w http.ResponseWriter, r *http.Request) {
	id := idMaker.NewId()
	idMsg := msgdef.SIdMsg{Op: msgdef.SIdOp, Id: id}
	if buf, err := json.Marshal(idMsg); err != nil {
		println(err.Error())
	} else {
		w.Write(buf)
	}
}

// Simple file server for serving up static content from the /html/ directory
// Also provides a simple id service for AJAX convenience
func main() {
	flag.Parse()
	logutil.ServerStarted("Example")
	pwd, err := os.Getwd()
	if err != nil {
		println(err.Error())
		return
	}
	locserver.StartTreeManager(10000, true)
	http.Handle("/loc", websocket.Handler(locserver.HandleLocationService))
	http.Handle("/msg", websocket.Handler(msgserver.HandleMessageService))
	http.HandleFunc("/id", idProvider)
	http.Handle("/", http.FileServer(http.Dir(pwd+"/html/")))
	portStr := fmt.Sprintf(":%d", *port)
	println(fmt.Sprintf("Listening on port %s", portStr))
	if err := http.ListenAndServe(portStr, nil); err != nil {
		println(err.Error())
	}
}
