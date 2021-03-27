package robot

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

var origin = "http://127.0.0.1:9989/"
var url = "ws://127.0.0.1:9989/event"
var WsCon *websocket.Conn
var wsLock sync.Mutex

func WebSocketClient2() {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("错误",err)
			}
		}
	}()

	var err error
	bk := false
	for {
		defer func() {
			if err := recover(); err != nil { //产生了panic异常
				if logger != nil {
					logger.Println("GetOnlineQQs异常:",err)
				}else {
					writeFile("exc.txt","onStart异常21\n")
				}
			}
		}()
		WsCon,err = websocket.Dial(url, "", origin);
		if err == nil {
			if logger != nil {
				logger.Println("ws链接成功")
			}
			go func() {
				err := recover()
				if err != nil {
					if logger != nil {
						logger.Println("GetOnlineQQs异常:",err)
					}else {
						writeFile("exc.txt","onStart异常21\n")
					}
				}
				GetOnlineQQs()
			}()
			break
		}else {
			wsLock.Lock()
			WsCon = nil
			wsLock.Unlock()
			time.Sleep(time.Second)
		}
	}
	go func() {
		defer func() {
			if err := recover(); err != nil { //产生了panic异常
				if logger != nil {
					logger.Println("GetOnlineQQs异常:",err)
				}else {
					writeFile("exc.txt","onStart异常22\r\n")
				}
			}
		}()
		for {
			if WsCon == nil {
				go func() {
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
					WebSocketClient2()
				}()
				//continue
				return
			}
			request := make([]byte, 2048)
			readLen, err := WsCon.Read(request)
			if readLen == 0 {
				wsLock.Lock()
				bk = true
				WsCon = nil
				wsLock.Unlock()
				go func() {
					defer func() {
						if err := recover(); err != nil { //产生了panic异常
							if logger != nil {
								logger.Println("GetOnlineQQs异常:",err)
							}else {
								writeFile("exc.txt","onStart异常98")
							}
							//WsCon = nil
						}
					}()
					WebSocketClient2()
				}()
				break;
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
			if err != nil {
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
				}()
				wsLock.Lock()
				bk = true
				WsCon = nil
				wsLock.Unlock()
				WebSocketClient2()
				break
			}
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
			}
			time.Sleep(time.Second * 30)
		}
	}()
}