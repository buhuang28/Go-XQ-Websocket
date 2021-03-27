package robot

func init() {
	createLog()
}

func onStart() {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("onStart异常:",err)
			}else {
				writeFile("exc.txt","onStart4异常\r\n")
			}
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil { //产生了panic异常
				if logger != nil {
					logger.Println("GetOnlineQQs异常:",err)
				}else {
					writeFile("exc.txt","onStart2异常\n")
				}
			}
		}()
		WebSocketClient2()
	}()
	go WebStart()
}

func onDisable() {

}

