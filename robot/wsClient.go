package robot

import (
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

var origin = "http://127.0.0.1:9989/"
var url = "ws://127.0.0.1:9989/event"
var WsCon *websocket.Conn
var wsLock sync.Mutex

func WebSocketClient2() {
	var err error
	bk := false
	for {
		WsCon,err = websocket.Dial(url, "", origin);
		if err == nil {
			break
		}else {
			wsLock.Lock()
			WsCon = nil
			wsLock.Unlock()
			time.Sleep(time.Second)
		}
	}
	request := make([]byte, 2048);
	go func() {
		for {
			if WsCon == nil {
				continue
			}
			readLen, err := WsCon.Read(request)
			if readLen == 0 {
				wsLock.Lock()
				bk = true
				WsCon = nil
				wsLock.Unlock()
				go WebSocketClient2()
				//条件变量同步通知
				break;
			} else {
				//处理websocket服务端发送过来的消息
				go processWebsocketMsg(request[:readLen])
			}
			if err != nil {
				wsLock.Lock()
				bk = true
				WsCon = nil
				wsLock.Unlock()
				go WebSocketClient2()
				break
			}
		}
	}()
	//这里不断向服务器那边传递在线QQ信息
	go func() {

		for  {
			if bk {
				return
			}
			if WsCon != nil {
				GetOnlineQQs()
			}else {
				return
			}
			time.Sleep(time.Second * 3)
		}
	}()
	if err := recover(); err != nil { //产生了panic异常
		if logger != nil {
			logger.Println("ws客户端异常:",err)
		}else {
			writeFile("exc.txt","wsclient异常")
		}
	}
}