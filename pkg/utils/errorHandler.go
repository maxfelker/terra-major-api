package utils

import (
	"log"
)

func ErrorHandler(e error) {
	if e != nil {
		log.Panic(e)
	}
}
