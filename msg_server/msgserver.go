package main

import (
	"errors"
	"fmt"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/msgwriter"
	"net/http"
	"websocket"
)

var idMap = simpleid.NewIdMap()

type user struct {
	id        string
	msgWriter *msgwriter.W
}

func newUser(ws *websocket.Conn) *user {
	return &user{msgWriter: msgwriter.New(ws)}
}

func readWS(ws *websocket.Conn) {
	var tId uint
	usr := newUser(ws)
	idMsg := msgdef.NewCIdMsg()
	if err := jsonutil.JSONCodec.Receive(ws, idMsg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	processReg(idMsg, usr)
	if err := idMap.Add(usr.id, usr); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	logutil.Registered(tId, usr.id)
	defer removeUser(&tId, usr.id)
	for {
		tId++
		cMsg := msgdef.NewCMsgMsg()
		if err := jsonutil.JSONCodec.Receive(ws, cMsg); err != nil {
			usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
			return
		}
		msg := cMsg.Msg.(*msgdef.CMsgMsg)
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user)
			msgMsg := &msgdef.SMsgMsg{From: usr.id, Content: msg.Content}
			forUser.msgWriter.WriteMsg(msgdef.NewServerMsg(msgdef.SMsgOp, msgMsg))
			logutil.Log(tId, usr.id, fmt.Sprintf("Content: '%s' send to: '%s'", msg.Content, msg.To))
		} else {
			usr.msgWriter.ErrorAndClose(tId, usr.id, fmt.Sprintf("User: %s is not present", msg.To))
			return
		}
	}
}

func processReg(clientMsg *msgdef.ClientMsg, usr *user) error {
	idMsg := clientMsg.Msg.(*msgdef.CIdMsg)
	switch clientMsg.Op {
	case msgdef.CAddOp:
		usr.id = idMsg.Id
		return nil
	}
	return errors.New("Incorrect op-code for id registration: " + string(clientMsg.Op))
}

func removeUser(tId *uint, uId string) {
	(*tId)++
	if idMap.Contains(uId) {
		idMap.Remove(uId)
		logutil.Deregistered(*tId, uId)
	} else {
		panic(fmt.Sprintf("User: %s\t Could not be removed from the message network"))
	}
}

func main() {
	logutil.ServerStarted("Message")
	http.Handle("/msg", websocket.Handler(readWS))
	http.ListenAndServe(":8003", nil)
}
