package main

import (
	"time"
	"finalproject123/server"
)

func main(){
	go server.Run()

	for {
		time.Sleep(5 * time.Second)
	}

}