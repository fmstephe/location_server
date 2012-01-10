package main

import (
	"os"
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

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		println(err.Error())
		return
	}
	http.HandleFunc("/id", idProvider)
	http.Handle("/", http.FileServer(http.Dir(pwd+"/html/")))
	if err := http.ListenAndServe(":8001", nil); err != nil {
		println(err.Error())
	}
}
