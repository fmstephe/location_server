package main

import (
	"io/ioutil"
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

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	println("index request")
	if file, err := ioutil.ReadFile("index.html"); err != nil {
		println(err.Error())
		return
	} else {
		w.Write(file)
	}
}

func main() {
	http.HandleFunc("/id", idProvider)
	http.HandleFunc("/example", exampleHandler)
	http.ListenAndServe(":8001", nil)
}
