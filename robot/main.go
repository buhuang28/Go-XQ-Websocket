package robot

func init() {
}

func onStart() {
	go createLog()
	go WebSocketClient2()
	go WebStart()
	if err := recover(); err != nil { //产生了panic异常
		logger.Println("onStart异常:",logger)
	}
}

func onDisable() {
}

