package robot

import (
	"encoding/json"
	"xq-go-sdk/core"
)

//先驱事件

//群聊
func OnGroupMessage(selfID,messageType,groupID,userID,messageNum,messageID,messageTime int64,message string)  {
	//需要的参数
	//xe.SelfID
	//xe.GroupID
	//xe.Time
	//xe.Message
	//撤回消息用的-----------
	//xe.MessageID
	//xe.MessageNum
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	xe := XEvent{
		SelfID:      selfID,		//接收消息的机器人QQ
		MessageType: messageType,
		GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
		UserID:      userID,		//消息来源ID
		Message:     message,		//消息内容
		MessageNum:  messageNum,	//用于消息撤回
		MessageID:   messageID,		//用于消息撤回
		Time:        messageTime,			//接收到消息的时间戳
	}
	marshal, _ := json.Marshal(xe)
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}

//私聊
func OnPrivateMessage(selfID,messageType,groupID,userID int64,message string)  {

	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	xe := XEvent{
		SelfID:      selfID,		//接收消息的机器人QQ
		MessageType: messageType,	//消息类型
		GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
		UserID:      userID,		//消息来源ID
		Message:     message,		//消息内容
	}
	marshal, _ := json.Marshal(xe)
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}

//群成员++
func OnGroupMemberIncrease(selfID,messageType,groupID,noticID int64)  {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	//NoticID //新增群员ID
	xe := XEvent{
		SelfID:      selfID,		//接收消息的机器人QQ
		MessageType: messageType,
		GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
		NoticID:      noticID,		//消息来源ID
	}
	nick := core.GetNick(selfID, noticID)
	xe.Name = nick
	marshal, _ := json.Marshal(xe)
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}

//回音信息处理---给儿子找爸爸
func OnEchoMessage(msg string,messageID,messageNum int64)  {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
		msgLock.Unlock()
	}()
	msgLock.Lock()
	if messageID != -1 && messageNum != -1 && messageID != 0 && messageNum != 0 {
		msgNumToID[messageNum] = messageID
	}
	logger.Println(msg,"\r\n的序号为:",messageNum,",消息ID为",messageID)
}

//成员被踢
func OnGroupMemberBeKick(xe *XEvent) {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	//xe.SelfID

}

//成员主动退群
func OnGroupMemberReduce(xe *XEvent)  {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	//xe.SelfID

}

//有人申请进群
func OnGroupMemberRequest(xe *XEvent)  {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	//xe.SelfID
	marshal, err := json.Marshal(xe)
	logger.Println("入群信息",string(marshal))
	if err != nil {
		logger.Println("申请入群序列化失败")
		return
	}
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}

//有人被邀请进群
func OnGroupMemberIntive(xe *XEvent) {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	marshal, err := json.Marshal(xe)
	logger.Println("被邀请入群信息",string(marshal))
	if err != nil {
		logger.Println("被邀请入群信息序列化失败")
		return
	}
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}

//有人邀请机器人进群
func OnIviteRobotToGroup(xe *XEvent){
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	marshal, err := json.Marshal(xe)
	logger.Println("有人邀请机器人进群信息",string(marshal))
	if err != nil {
		logger.Println("有人邀请机器人进群序列化失败")
		return
	}
	if 	WsCon != nil {
		_, _ = WsCon.Write(marshal)
	}
}