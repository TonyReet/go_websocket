package main

import "mc/src/service/wSocket"

func main() {
	ch := make(chan bool)

	go wSocket.InitWSocket()

	<-ch
}
