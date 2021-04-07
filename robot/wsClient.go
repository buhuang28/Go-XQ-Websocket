package robot

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"sync"
	"time"
	"xq-go-sdk/core"
)

var origin = "http://127.0.0.1:9989/"
var url = "ws://127.0.0.1:9989/event"
var WsCon *websocket.Conn = nil
var wsLock sync.Mutex


func WebSocketClient2() {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("错误",err)
			}
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil { //产生了panic异常
				if logger != nil {
					logger.Println("GetOnlineQQs异常:",err)
				}else {
					writeFile("exc.txt","onStart异常99\n")
				}
				go WebSocketClient2()
			}
		}()
		for {
			defer func() {
				if err := recover(); err != nil { //产生了panic异常
					if logger != nil {
						logger.Println("GetOnlineQQs异常:",err)
					}else {
						writeFile("exc.txt","onStart异常99\n")
					}
					//WsCon = nil
				}
			}()
			if WsCon == nil {
				go func() {
					Websocket_Client_Dail()
				}()
				time.Sleep(time.Second*3)
				continue
			}
			request := make([]byte, 2048)
			readLen, err := WsCon.Read(request)
			if readLen == 0 ||err != nil {
				wsLock.Lock()
				WsCon = nil
				wsLock.Unlock()
				go func() {
					Websocket_Client_Dail()
				}()
				time.Sleep(time.Second*3)
				continue;
			} else {
				//处理websocket服务端发送过来的消息
				var req []byte
				req = request[:readLen]
				go func() {
					defer func() {
						if err := recover(); err != nil { //产生了panic异常
							if logger != nil {
								logger.Println("GetOnlineQQs异常:",err)
							}else {
								writeFile("exc.txt","onStart异常")
							}
						}
					}()
					processWebsocketMsg(req)
				}()
			}
			//if err != nil {
			//	defer func() {
			//		if err := recover(); err != nil { //产生了panic异常
			//			if logger != nil {
			//				logger.Println("GetOnlineQQs异常:",err)
			//			}else {
			//				writeFile("exc.txt","onStart异常")
			//			}
			//		}
			//	}()
			//	wsLock.Lock()
			//	WsCon = nil
			//	wsLock.Unlock()
			//	time.Sleep(time.Second*3)
			//	go func() {
			//		Websocket_Client_Dail()
			//	}()
			//	continue
			//}
		}
	}()
	//这里不断向服务器那边传递在线QQ信息   -----废弃，改为websocket心跳包
	//暂定websocket心跳
	go func() {
		defer func() {
			if err := recover(); err != nil { //产生了panic异常
				if logger != nil {
					logger.Println("GetOnlineQQs异常:",err)
				}else {
					writeFile("exc.txt","onStart异常")
				}
			}
			WsCon = nil
			go Websocket_Client_Dail()
		}()

		go func() {
			for {
				err := recover()
				if err != nil {
					logger.Println(err)
				}
				core.FuckXQ()
				time.Sleep(time.Second * 5)
			}
		}()

		for  {
			if WsCon != nil {
				//GetOnlineQQs()
				var xe ApiCallBack
				xe.Robot = 1
				xe.Type = HEART_BEATS
				marshal, err2 := json.Marshal(xe)
				if err2 == nil && WsCon != nil {
					WsCon.Write(marshal)
				}
			}else {
				go Websocket_Client_Dail()
			}
			time.Sleep(time.Second * 20)
		}
	}()
}

func Websocket_Client_Dail() {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("GetOnlineQQs异常:",err)
			}else {
				writeFile("exc.txt","onStart异常21\n")
			}
		}
	}()
	var err error
	wsLock.Lock()
	WsCon,err = websocket.Dial(url, "", origin);
	logger.Println("Dail错误",err)
	wsLock.Unlock()
	if WsCon != nil && err == nil {
		time.Sleep(time.Second * 3)
		GetOnlineQQs()
		if logger != nil {
			logger.Println("ws链接成功")
		}
		return
	}else {
		wsLock.Lock()
		WsCon = nil
		wsLock.Unlock()
		//time.Sleep(time.Second)
		//go Websocket_Client_Dail()
	}
}