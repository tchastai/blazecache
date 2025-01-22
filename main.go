package main

import (
	"blazecache/server"
	"log"
)

func main() {
	logChan := make(chan string)
	agentAdress := []string{":8070", ":8071", ":8072"}
	for k, v := range agentAdress {
		go func(k int, v string) {
			a := server.New(k, v, logChan, agentAdress)
			a.Start()
		}(k, v)
	}

	for msg := range logChan {
		log.Println(msg)
	}
}
