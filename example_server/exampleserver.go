package main

import (
	"encoding/json"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/msgdef"
	"net/http"
	"os"
	"os/exec"
)

var idMaker = simpleid.NewIdMaker()

func idProvider(w http.ResponseWriter, r *http.Request) {
	id := idMaker.NewId()
	idMsg := msgdef.SIdMsg{Op: msgdef.SIdOp, Id: id}
	if buf, err := json.MarshalForHTML(idMsg); err != nil {
		println(err.Error())
	} else {
		w.Write(buf)
	}
}

func restart(w http.ResponseWriter, r *http.Request) {
	os.Chdir("scripts")
	cmd := exec.Command("./update_servers.sh")
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
}

func main() {
	logutil.ServerStarted("Example")
	pwd, err := os.Getwd()
	if err != nil {
		println(err.Error())
		return
	}
	http.HandleFunc("/id", idProvider)
	http.HandleFunc("/restart", restart)
	http.Handle("/", http.FileServer(http.Dir(pwd+"/html/")))
	if err := http.ListenAndServe(":8001", nil); err != nil {
		println(err.Error())
	}
}
