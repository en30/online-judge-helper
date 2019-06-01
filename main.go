package main

import (
	"log"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	sc := make(chan *Submission)
	server := newServer(config)
	go watch(sc, config)
	go server.launch()

	for {
		submission := <-sc
		server.sendSubmission(submission)
	}
}
