package main

import (
	"errors"
	"fmt"
	"net/http"
	"websocket"
	"encoding/json"
	"location_server/user"
	"location_server/msgutil/msgdef"
	"location_server/msgutil/jsonutil"
	"github.com/fmstephe/simpleid"
)

var idMap = simpleid.NewIdMap()

func readWS(ws *websocket.Conn) {
	idMsg := msgdef.NewCIdMsg()
	if err := jsonutil.JSONCodec.Receive(ws, idMsg); err != nil {
		fmt.Printf("Msg Srvr: Connection Terminated with %s\n", err.Error())
		return
	}
	usr := user.New()
	processReg(idMsg, usr)
	if err := idMap.Add(usr.Id, usr); err != nil {
		fmt.Println("Msg Srvr: Connection Terminated with ", err.Error())
		return
	} else {
		fmt.Printf("Msg Srvr: User: %s\tAdded successfully to the message network\n", usr.Id)
	}
	defer idMap.Remove(usr.Id)
	go writeWS(ws, usr)
	for {
		cMsg := msgdef.NewCMsgMsg()
		if err := jsonutil.JSONCodec.Receive(ws, cMsg); err != nil {
			fmt.Printf("Msg Srvr: User: %s\tConnection Terminated with %s\n", usr.Id, err.Error())
			return
		}
		msg := cMsg.Msg.(*msgdef.CMsgMsg)
		if idMap.Contains(msg.To) {
			forUser := idMap.Get(msg.To).(*user.U)
			msgMsg := &msgdef.SMsgMsg{From: usr.Id, Content: msg.Content}
			forUser.WriteMsg(msgdef.NewServerMsg(msgdef.SMsgOp, msgMsg))
		} else {
			fmt.Printf("Msg Srvr: User: %s\t Is not present\n", msg.To)
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

func writeWS(ws *websocket.Conn, usr *user.U) {
	defer closeWS(ws, usr)
	for {
		msg := usr.ReceiveMsg()
		buf, err := json.MarshalForHTML(msg)
		if err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
			return
		}
		fmt.Printf("User: %s \tServer Message: %s\n", usr.Id, string(buf))
		if _, err = ws.Write(buf); err != nil {
			fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
			return
		}
	}
}

func closeWS(ws *websocket.Conn, usr *user.U) {
	if err := ws.Close(); err != nil {
		fmt.Printf("User: %s \tError: %s\n", usr.Id, err.Error())
	}
}

func main() {
	http.Handle("/msg", websocket.Handler(readWS))
	http.ListenAndServe(":8003", nil)
}
