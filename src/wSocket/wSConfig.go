package wSocket

import (
	"flag"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

var (
	//服务器地址
	addr = flag.String("addr", "localhost:8081", "service address")

	//心跳最大时间
	heartbeatTimeOut = flag.Int64("heartbeatTimeOut", 20, "heartbeat TimeOut")

	//超过心跳的次数
	heartbeatTimeOutCount = flag.Int64("timeOutCount", 3, "heartbeat TimeOut Count")

	upgrader = websocket.Upgrader{
		EnableCompression: false,
		ReadBufferSize:    256,
		WriteBufferSize:   256,

		//不用检查Origin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

		//设置握手的时间
		HandshakeTimeout: time.Duration(2) * time.Minute,
	}
)
