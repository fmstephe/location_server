package logutil

import (
	"log"
)

func Log(tId uint, uId, msg string) {
	log.Printf("%d\t%s\t%s",tId, uId, msg)
}

func LogFree(msg string) {
	log.Print(msg)
}

func Connected() {
	log.Printf("Connection Established")
}

func Registered(tId uint, uId string) {
	Log(tId, uId, "User Registered")
}

func Deregistered(tId uint, uId string) {
	Log(tId, uId, "User Deregistered")
}
