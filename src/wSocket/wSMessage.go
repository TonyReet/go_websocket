package wSocket

import (
	"io"
	"time"
	"fmt"
	"log"
	"mc/src/utils"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	//保存conn的key字符串
	keyWord string

	msg interface{}
}

var (
	//写消息队列
	writeMessages = make(chan *WsMessage)
)

func (wS *WebSocket)handleNewMessage(conn *websocket.Conn){
	wS.lock.Lock()

	addrStr := fmt.Sprintf("%p", &conn)
	keyWord := conn.RemoteAddr().String() + conn.RemoteAddr().Network() + addrStr

	//如果存在keyWord
	if len(keyWord) != 0 {

		//删除旧的conn,因为旧的已经失效
		_, ok := wS.conns[keyWord]
		if ok {
			delete(wS.conns, keyWord)
		}

		log.Println("新的连接地址:", keyWord)
		wS.conns[keyWord] = conn
	}
	wS.lock.Unlock()

	//读操作
	wS.readLoop(conn, keyWord)
}

func (wS *WebSocket) readLoop(conn *websocket.Conn, keyWord string) {
	//首次进入设置超时时间，SetReadDeadline设置的时间超过以后就不再读取
	conn.SetReadDeadline(time.Now().Add(timeOut * time.Second))

	for {

		msgType, msg, err := wS.readMessage(conn)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				log.Println("webSocket 已经关闭:", err)
			} else {
				log.Println("websocket 读取出错", err)
			}

			conn.Close()
			break
		}

		switch msgType {
		case websocket.TextMessage: //文本消息就重设读取时间

			wS.WSWriteTextMessages(msg)

			//重置超时时间
			conn.SetReadDeadline(time.Now().Add(timeOut * time.Second))
		default:
			log.Println("不支持的消息类型:", msgType)
		}
	}

	//发送关闭消息
	conn.WriteMessage(websocket.CloseMessage, []byte{})
}

/*读取消息*/
func (wS *WebSocket) readMessage(conn *websocket.Conn) (msgType int, p []byte, err error) {
	msgType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("读取消息出错:", err)


		return msgType, message, err
	}

	log.Printf("读取的消息: %s", message)

	return msgType, message, err
}

/*发送文本消息*/
func (wS *WebSocket) writeTextMessage(keyWord string, msg string) error {

	conn := wS.conns[keyWord]

	return wS.writeMessage(conn, websocket.TextMessage, msg)
}

/*发送消息的统一处理*/
func (wS *WebSocket) writeMessage(conn *websocket.Conn, msgType int, msg string) error {

	messageFinal := []byte(msg)
	err := conn.WriteMessage(msgType, messageFinal)
	if err != nil {
		log.Println("发送消息出错:", err)
		return err
	}

	log.Printf("发送消息: %s,成功", msg)
	return nil

}

/*处理写入的消息*/
func (wS *WebSocket) HandleWriteMessages() {
	for {
		wsMsg := <-writeMessages
		keyWord := wsMsg.keyWord
		msg := wsMsg.msg

		switch msg.(type) {
		case string:
			if message, ok := msg.(string); ok {
				if len(keyWord) != 0 {
					wS.writeTextMessage(keyWord, message)
				} else {
					log.Println("<写>没有对应的keyWord信息")
				}
			} else {
				log.Println("<写>类型转换失败:")
			}
		case interface{}:
			log.Println("interface")
		default:
			log.Println("<写>不支持的类型")
		}
	}
}

/*写入*/
func (wS *WebSocket) WSWriteTextMessages(wsMessage *WsMessage) {
	writeMessages <- wsMessage
}