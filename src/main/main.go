package main

import "go_websocket/src/wSocket"

func main() {
	ch := make(chan bool)

	go wSocket.InitWSocket()

	<-ch
}
