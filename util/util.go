package util

import (
	"log"
)

var VerboseLogging = true

// Should always log erros
func LogErr(err error) {
	log.Printf("[ERROR] %v\n", err.Error())
}

func LogInfo(msg string) {
	if VerboseLogging {
		log.Println("[INFO] ", msg)
	}
}

func Log(msg string) {
	if VerboseLogging {
		log.Println(msg)
	}
}
