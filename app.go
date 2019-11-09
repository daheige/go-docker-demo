package main

import (
	"log"
	"time"
)

func main() {
	log.Println(111)
	log.Println("hello")

	for {
		log.Println(111)
		time.Sleep(1 * time.Second)
	}
}
