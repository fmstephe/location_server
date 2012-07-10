package msgserver

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"github.com/fmstephe/location_server/logutil"
	"github.com/fmstephe/location_server/msgutil/jsonutil"
	"github.com/fmstephe/location_server/msgutil/msgdef"
	"github.com/fmstephe/location_server/user"
	"github.com/fmstephe/simpleid"
)

var idMap = simpleid.NewIdMap()

func HandleMessageService(ws *websocket.Conn) {
	var tId uint
	usr := user.New(ws)
	idMsg := &msgdef.CIdMsg{}
	procReg := processReg(tId, idMsg, usr)
	if err := jsonutil.UnmarshalAndProcess(tId, usr.Id, ws, idMsg, procReg); err != nil {
		usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
		return
	}
	defer removeUser(&tId, usr.Id)
	for {
		tId++
		msg := &msgdef.CMsgMsg{}
		procMsg := processMsg(tId, msg, usr)
		if err := jsonutil.UnmarshalAndProcess(tId, usr.Id, ws, msg, procMsg); err != nil {
			usr.MsgWriter.ErrorAndClose(tId, usr.Id, err.Error())
			return
		}
	}
}

func processReg(tId uint, idMsg *msgdef.CIdMsg, usr *user.U) func() error {
	return func() error {
		if idMsg.Op != msgdef.CAddOp {
			return errors.New("Incorrect op-code for id registration: " + string(idMsg.Op))
		}
		if err := idMsg.Validate(); err != nil {
			return err
		}
		usr.Id = idMsg.Id
		if err := idMap.Add(usr.Id, usr); err != nil {
			return err
		}
		logutil.Registered(tId, usr.Id)
		return nil
	}
}

func processMsg(tId uint, msg *msgdef.CMsgMsg, usr *user.U) func() error {
	return func() error {
		if msg.Op != msgdef.CMsgOp {
			return errors.New("Incorrect op-code for msg: " + string(msg.Op))
		}
		if err := msg.Validate(); err != nil {
			return err
		}
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user.U)
			safeContent := jsonutil.SanitiseJSON(msg.Content)
			msgMsg := &msgdef.SMsgMsg{Op: msgdef.SMsgOp, From: usr.Id, Content: safeContent}
			sMsg := &msgdef.ServerMsg{Msg: msgMsg, TId: tId, UId: usr.Id}
			forUser.MsgWriter.WriteMsg(sMsg)
			logutil.Log(tId, usr.Id, fmt.Sprintf("Content: '%s' sent to: '%s'", msg.Content, msg.To))
		} else {
			nuMsg := &msgdef.SMsgMsg{Op: msgdef.SNotUserOp, From: msg.To, Content: fmt.Sprintf("User: %s was not found", msg.To)}
			sMsg := &msgdef.ServerMsg{Msg: nuMsg, TId: tId, UId: usr.Id}
			usr.MsgWriter.WriteMsg(sMsg)
		}
		return nil
	}
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
