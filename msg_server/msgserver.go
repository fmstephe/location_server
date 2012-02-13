package main

import (
	"encoding/json"
	"fmt"
	"github.com/fmstephe/simpleid"
	"location_server/msgdef"
	"net/http"
	"websocket"
)

var idMap = simpleid.NewIdMap()

type user struct {
	id      string
	msgChan chan *msgdef.SMsgMsg
}

func readWS(ws *websocket.Conn) {
	buf := make([]byte, 2024)
	idMsg := new(msgdef.CIdMsg)
	if err := unmarshal(buf, idMsg, ws); err != nil {
		fmt.Printf("Connection Terminated with %s\n", err.Error())
		return
	}
	usr := &user{id: idMsg.Id, msgChan: make(chan *msgdef.SMsgMsg, 16)}
	if err := idMap.Add(idMsg.Id, usr); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Printf("User: %s\tAdded successfully to the message network\n", usr.id, )
	}
	defer idMap.Remove(idMsg.Id)
	go writeWS(ws, usr)
	for {
		msg := new(msgdef.CMsgMsg)
		if err := unmarshal(buf, msg, ws); err != nil {
			fmt.Printf("User: %s\tConnection Terminated with %s\n", usr.id, err.Error())
			return
		}
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user)
			forUser.msgChan <- &msgdef.SMsgMsg{Op: msgdef.SMsgOp, From: usr.id, Msg: msg.Msg}
		} else {
			fmt.Printf("User: %s\t Is not present\n", msg.To)
			return
		}
	}
}

func removeUser(id string) {
	if idMap.Contains(id) {
		idMap.Remove(id)
		fmt.Printf("User: %s\t Successfully removed from the message network\n")
	} else {
		panic(fmt.Sprintf("User: %s\t Could not be removed from the message network"))
	}
}

func unmarshal(buf []byte, msg interface{}, ws *websocket.Conn) error {
	err := websocket.JSON.Receive(ws, msg)
	if err != nil {
		return err
	}
	return nil
}

func writeWS(ws *websocket.Conn, usr *user) {
	defer closeWS(ws, usr)
	for {
		msg := <-usr.msgChan
		buf, err := json.MarshalForHTML(msg)
		if err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.id, err.Error())
			return
		}
		fmt.Printf("User: %s \tServer Message: %s\n", usr.id, string(buf))
		if _, err = ws.Write(buf); err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.id, err.Error())
			return
		}
	}
}

func closeWS(ws *websocket.Conn, usr *user) {
	if err := ws.Close(); err != nil {
		fmt.Printf("User: %s \tError: %s\n", usr.id, err.Error())
	}
}

func main() {
	http.Handle("/msg", websocket.Handler(readWS))
	http.ListenAndServe(":8003", nil)
}
