package util

import (
	"log"
)

func LogErr(err error) {
	log.Printf("[ERROR] %v\n", err.Error())
}
