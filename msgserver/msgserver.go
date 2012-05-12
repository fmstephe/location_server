package msgserver

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/msgwriter"
)

var idMap = simpleid.NewIdMap()

type user struct {
	id        string
	msgWriter *msgwriter.W
}

func newUser(ws *websocket.Conn) *user {
	return &user{msgWriter: msgwriter.New(ws)}
}

func HandleMessageService(ws *websocket.Conn) {
	var tId uint
	usr := newUser(ws)
	idMsg := &msgdef.CIdMsg{}
	if err := jsonutil.JSONCodec.Receive(ws, idMsg); err != nil {
		usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
		return
	}
	if err := idMsg.Validate(); err != nil {
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
		msg := &msgdef.CMsgMsg{}
		if err := jsonutil.JSONCodec.Receive(ws, msg); err != nil {
			usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
			return
		}
		if err := msg.Validate(); err != nil {
			usr.msgWriter.ErrorAndClose(tId, usr.id, err.Error())
			return
		}
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user)
			msgMsg := &msgdef.SMsgMsg{Op: msgdef.SMsgOp, From: usr.id, Content: msg.Content}
			sMsg := &msgdef.ServerMsg{Msg: msgMsg, TId: tId, UId: usr.id}
			forUser.msgWriter.WriteMsg(sMsg)
			logutil.Log(tId, usr.id, fmt.Sprintf("Content: '%s' sent to: '%s'", msg.Content, msg.To))
		} else {
			nuMsg := &msgdef.SMsgMsg{Op: msgdef.SNotUserOp, From: msg.To, Content: fmt.Sprintf("User: %s was not found", msg.To)}
			sMsg := &msgdef.ServerMsg{Msg: nuMsg, TId: tId, UId: usr.id}
			usr.msgWriter.WriteMsg(sMsg)
		}
	}
}

func processReg(idMsg *msgdef.CIdMsg, usr *user) error {
	switch idMsg.Op {
	case msgdef.CAddOp:
		usr.id = idMsg.Id
		return nil
	}
	return errors.New("Incorrect op-code for id registration: " + string(idMsg.Op))
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
