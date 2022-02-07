package main

import "log"

func logIfError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
