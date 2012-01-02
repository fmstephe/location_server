package main

import (
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
	http.HandleFunc("/", idProvider)
	http.ListenAndServe(":8001", nil)
}
