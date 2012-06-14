package msgserver

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"github.com/fmstephe/simpleid"
	"location_server/logutil"
	"location_server/msgutil/jsonutil"
	"location_server/msgutil/msgdef"
	"location_server/user"
)

var idMap = simpleid.NewIdMap()

func HandleMessageService(ws *websocket.Conn) {
	var tId uint
	usr := user.New(ws)
	idMsg := &msgdef.CIdMsg{}
	if err := jsonutil.JSONCodec.Receive(ws, idMsg); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	if err := idMsg.Validate(); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	processReg(idMsg, usr)
	if err := idMap.Add(usr.Id, usr); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	logutil.Registered(tId, usr.Id)
	defer removeUser(&tId, usr.Id)
	for {
		tId++
		msg := &msgdef.CMsgMsg{}
		if err := jsonutil.JSONCodec.Receive(ws, msg); err != nil {
			usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
			return
		}
		if err := msg.Validate(); err != nil {
			usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
			return
		}
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user.U)
			msgMsg := &msgdef.SMsgMsg{Op: msgdef.SMsgOp, From: usr.Id, Content: msg.Content}
			sMsg := &msgdef.ServerMsg{Msg: msgMsg, TId: tId, UId: usr.Id}
			forUser.MsgWriter.WriteMsg(sMsg)
			logutil.Log(tId, usr.Id, fmt.Sprintf("Content: '%s' sent to: '%s'", msg.Content, msg.To))
		} else {
			nuMsg := &msgdef.SMsgMsg{Op: msgdef.SNotUserOp, From: msg.To, Content: fmt.Sprintf("User: %s was not found", msg.To)}
			sMsg := &msgdef.ServerMsg{Msg: nuMsg, TId: tId, UId: usr.Id}
			usr.MsgWriter.WriteMsg(sMsg)
		}
	}
}

func processReg(idMsg *msgdef.CIdMsg, usr *user.U) error {
	switch idMsg.Op {
	case msgdef.CAddOp:
		usr.Id = idMsg.Id
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
