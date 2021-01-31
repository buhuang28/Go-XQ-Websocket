package robot
//
//import "time"
//
//package robot
//
//import (
//"golang.org/x/net/websocket"
//"sync"
//"time"
//)
//
//var origin = "http://127.0.0.1:9989/"
//var url = "ws://127.0.0.1:9989/event"
//var WsCon *websocket.Conn
//var wsLock sync.Mutex //群Map锁
//
//
//func WebSocketClient() {
//	var err error
//	bk := false
//	for {
//		wsLock.Lock()
//		WsCon,err = websocket.Dial(url, "", origin);
//		wsLock.Unlock()
//		if err == nil {
//			break
//		}else {
//			wsLock.Lock()
//			WsCon = nil
//			wsLock.Unlock()
//			time.Sleep(time.Second)
//			go WebSocketClient()
//			return
//		}
//	}
//	request := make([]byte, 2048);
//	go func() {
//		for {
//			if WsCon == nil {
//				go WebSocketClient()
//				bk = true
//				return
//			}
//			readLen, err := WsCon.Read(request)
//			if readLen == 0 {
//				go WebSocketClient()
//				bk = true
//				continue
//			} else {
//				//处理websocket服务端发送过来的消息
//				processWebsocketMsg(request[:readLen])
//
//			}
//			if err != nil {
//				bk = true
//				wsLock.Lock()
//				WsCon = nil
//				wsLock.Unlock()
//				go WebSocketClient()
//				return
//			}
//		}
//	}()
//	//这里不断向服务器那边传递在线QQ信息
//	go func() {
//		for  {
//			if bk || WsCon == nil {
//				return
//			}
//			GetOnlineQQs()
//			time.Sleep(time.Second * 3)
//		}
//	}()
//}
//
//
