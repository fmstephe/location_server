package main

import (
	"os"
	"os/exec"
	"encoding/json"
	"github.com/fmstephe/simpleid"
	"location_server/msgdef"
	"net/http"
)

var idMaker = simpleid.NewIdMaker()

func idProvider(w http.ResponseWriter, r *http.Request) {
	id := idMaker.NewId()
	idMsg := msgdef.NewSIdMsg(id)
	if buf, err := json.MarshalForHTML(idMsg); err != nil {
		println(err.Error())
	} else {
		w.Write(buf)
	}
}

func restart(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("./scripts/update_servers.sh")
	err := cmd.Run()
	if err != nil {
		println(err.Error())
	}
}

func main() {
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
