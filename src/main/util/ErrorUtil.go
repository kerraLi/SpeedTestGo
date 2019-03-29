package util

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func FailOnErrorNoExit(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
