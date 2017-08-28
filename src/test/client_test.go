package test

import (
	"flag"
	"log"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"github.com/gorilla/websocket"

	"testing"
)

var clientAddr = flag.String("addr", "localhost:8081", "http service address")
var sendMulti = flag.Bool("send multi message", true, "Whether send multi message")

func Test_send(t *testing.T) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for index := 0; index < 100; index++ {
		go newCoon(index)
	}

	//阻塞，保持程序不终止
	ch := make(chan bool)
	log.Print("all message send finish")
	<-ch
}

func newCoon(index int) {
	u := url.URL{Scheme: "ws", Host: *clientAddr, Path: "/ws"}

	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("is%dclient,recv: %s", index, message)
		}
	}()

	if *sendMulti == true {
		for sendIndex := 0; sendIndex < 10; sendIndex++ {
			message := fmt.Sprintln("is", index, "client send", sendIndex, "message")
			err2 := c.WriteMessage(websocket.TextMessage, []byte(message))

			if err2 != nil {
				log.Println("send fail:", err2)
				return
			}
		}
	} else {
		message := fmt.Sprintln("is", index, "client")
		err2 := c.WriteMessage(websocket.TextMessage, []byte(message))

		if err2 != nil {
			log.Println("send fail:", err2)
			return
		}

	}

	//阻塞，保持链接不断开
	ch := make(chan bool)
	<-ch
}
