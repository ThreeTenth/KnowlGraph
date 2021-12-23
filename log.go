package main

import (
	"log"
)

func info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func debug(format string, v ...interface{}) {
	if config.Debug {
		log.Printf(format, v...)
	}
}
