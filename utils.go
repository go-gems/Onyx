package Onyx

import (
	"log"
)

var Debug = false

func debugPrintln(i ...interface{}) {
	if Debug {
		log.Println(i...)
	}
}
func debugPrintf(format string, i ...interface{}) {
	if Debug {
		log.Printf(format, i...)
	}
}
