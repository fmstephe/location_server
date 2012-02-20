package main

import (
	"errors"
	"fmt"
	"net/http"
	"websocket"
	"encoding/json"
	"location_server/msgutil/msgwriter"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/jsonutil"
	"github.com/fmstephe/simpleid"
)

var idMap = simpleid.NewIdMap()

type user struct {
	id string
	msgWriter *msgWriter.W
}

func newUser(ws *websocket.Conn) {
	return &user{msgWriter: msgwriter.New(ws)}
}

func readWS(ws *websocket.Conn) {
	var tId uint
	usr := newUser(ws)
	idMsg := msgdef.NewCIdMsg()
	if err := jsonutil.JSONCodec.Receive(ws, idMsg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	processReg(idMsg, usr)
	if err := idMap.Add(usr.Id, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, err.Error())
		return
	}
	fmt.Printf("Msg Srvr: User: %s\tAdded successfully to the message network\n", usr.Id)
	defer idMap.Remove(usr.Id)
	for {
		tId++
		cMsg := msgdef.NewCMsgMsg()
		if err := jsonutil.JSONCodec.Receive(ws, cMsg); err != nil {
			usr.msgWriter.ErrorAndClose(tId, err.Error())
			return
		}
		msg := cMsg.Msg.(*msgdef.CMsgMsg)
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user.U)
			msgMsg := &msgdef.SMsgMsg{From: usr.Id, Content: msg.Content}
			forUser.WriteMsg(msgdef.NewServerMsg(msgdef.SMsgOp, msgMsg))
		} else {
			usr.msgWriter.ErrorAndClose(tId, fmt.Sprintf("User: %s is not present",msg.To)
			return
		}
	}
}

func processReg(clientMsg *msgdef.ClientMsg, usr *user.U) error {
	idMsg := clientMsg.Msg.(*msgdef.CIdMsg)
	switch clientMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
		return nil
	}
	return errors.New("Incorrect op-code for id registration: "+string(clientMsg.Op))
}

func removeUser(id string) {
	if idMap.Contains(id) {
		idMap.Remove(id)
		fmt.Printf("User: %s\t Successfully removed from the message network\n")
	} else {
		panic(fmt.Sprintf("User: %s\t Could not be removed from the message network"))
	}
}

func main() {
	http.Handle("/msg", websocket.Handler(readWS))
	http.ListenAndServe(":8003", nil)
}
