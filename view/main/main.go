package main

import (
	"chatroom/client"
	"chatroom/view"
	"flag"
	"log"
)

func main() {
	address := flag.String("server", "localhost:3333", "Which server to connect to")

	flag.Parse()

	//製造一個新的client端接server
	client := client.NewClient()
	//透過這裡連到server
	err := client.Dial(*address)

	if err != nil {
		log.Fatal(err)
	}

	//defer https://www.evanlin.com/golang-know-using-defer/
	defer client.Close()

	// start the client to listen for incoming message
	go client.Start()

	view.UI(client)
}
