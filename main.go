package main

import (
	"log"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	go watch(config)
	go newServer(config).launch()
	<-done
}
