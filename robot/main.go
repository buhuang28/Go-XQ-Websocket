package robot

func init() {
}

func onStart() {
	go createLog()
	go WebSocketClient2()
	go WebStart()
	if err := recover(); err != nil { //产生了panic异常
		if logger != nil {
			logger.Println("onStart异常:",err)
		}else {
			writeFile("exc,txt","onStart异常")
		}
	}
}

func onDisable() {
}

