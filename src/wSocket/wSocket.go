package wSocket

import (
	"flag"
	"log"
	"time"
	"sync"
	"net/http"
	"mc/src/utils"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	conns map[string]*websocket.Conn

	lock sync.Mutex
}

var (
	webSocket *WebSocket
	once      sync.Once

	//超时时间*超时次数
	timeOut = time.Duration(*heartbeatTimeOut) * time.Duration(*heartbeatTimeOutCount)
)

/*单例*/
func GetInstance() *WebSocket {
	once.Do(func() {
		webSocket = &WebSocket{}
	})
	return webSocket
}

func (wS *WebSocket) socket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not support", http.StatusMethodNotAllowed)
		return
	}

	//将http升级为websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	//处理新消息
	wS.handleNewMessage(conn)
}

func InitWSocket() {
	wS := GetInstance()

	flag.Parse()

	//输出时间+输出路径+文件名+行号格式
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("-------------------webSocket started-------------------")

	wS.conns = make(map[string]*websocket.Conn)
	http.HandleFunc("/", wS.socket)

	//写操作
	go wS.HandleWriteMessages()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("webSocket start fail")
	}
}
