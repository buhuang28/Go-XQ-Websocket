package robot

import (
	"strconv"
	"xq-go-sdk/core"
)

var AppInfoJson string

func init() {
	core.Create = XQCreate
	core.Event = XQEvent
	core.DestroyPlugin = XQDestroyPlugin
	core.SetUp = XQSetUp
}

func XQCreate(version string) string {
	return AppInfoJson
}

var (
	msgMap = make(map[string]int64)
)

type XEvent struct {
	SelfID      int64 `json:"self_id"`
	MessageType int64 `json:"message_type"`
	SubType     int64 `json:"sub_type"`
	GroupID     int64 `json:"group_id"`
	UserID      int64 `json:"user_id"`
	NoticID     int64 `json:"notic_id"`
	Message     string `json:"message"`
	MessageNum  int64 `json:"message_num"`
	MessageID   int64 `json:"message_id"`
	RawMessage  string `json:"raw_message"`
	Time        int64 `json:"time"`
	Ret         int64 `json:"ret"`
	CqID        int64 `json:"cq_id"`
	Name		string `json:"name"`
}

func XQEvent(selfID int64, messageType int64, subType int64, groupID int64, userID int64, noticID int64, message string, messageNum int64, messageID int64, rawMessage string, time int64, ret int64) int64 {
	switch messageType {
	//插件启动事件
	case 12001:
		go onStart()
	//插件关闭事件
	case 12002:
	// 消息事件下
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		go func() {
			OnPrivateMessage(selfID,messageType,groupID,userID,message)
		}()
	// 2：群聊信息
	case 2, 3:
		go func() {
			defer func() {
				if err := recover(); err != nil { //产生了panic异常
					logger.Println(err)
				}
			}()
			mesSid := strconv.FormatInt(messageID,10)+ strconv.FormatInt(selfID,10)
			msgLock.Lock()
			num,ok := msgMap[mesSid]
			msgLock.Unlock()
			if ok && num == messageNum {
				return
			}else {
				msgLock.Lock()
				msgMap[mesSid] = messageNum
				msgLock.Unlock()
				OnGroupMessage(selfID,messageType,groupID,userID,messageNum,messageID,time,message)
				if len(msgMap) > 10500 {
					msgLock.Lock()
					msgMap = make(map[string]int64)
					msgLock.Unlock()
				}
			}
		}()
	// 10：回音信息---发送消息后的回音
	case 10:
		go OnEchoMessage(message,messageID,messageNum)
	// 通知事件
	// 群文件接收
	case 218:
	// 管理员变动 210为有人升为管理 211为有人被取消管理
	case 210:
	case 211:
	// 群成员减少 201为主动退群 202为被踢
	case 201:
	case 202:
	// 群成员增加
	case 212:
		go OnGroupMemberIncrease(selfID,messageType,groupID,noticID)
	// 群禁言 203为禁言 204为解禁
	case 203:
	case 204:
	// new
	// 好友添加 100 为单向 102 为标准
	case 100, 102:

	// 群消息撤回 subType 2
	// 好友消息撤回 subType 1
	case 9:

	// 群内戳一戳

	// 群红包运气王

	// 群成员荣誉变更

	// 请求事件
	// 加好友请求
	case 101:
	// 加群请求／邀请 213请求入群  214我被邀请加入某群  215某人被邀请加入群
	case 213:
		go func() {
			nick := core.GetNick(selfID,userID)
			xe := XEvent {
				SelfID:      selfID,		//接收消息的机器人QQ
				MessageType: messageType,	//消息类型
				GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
				UserID:      userID,		//消息来源ID
				NoticID:     noticID,		//触发对象、被动
				Message:     message,		//消息内容
				RawMessage:  rawMessage,	//原始信息
				Name: 		 nick,
			}
			OnGroupMemberRequest(&xe)
		}()
	case 214:
		go func() {
			nick := core.GetNick(selfID,noticID)
			xe := XEvent {
				SelfID:      selfID,		//接收消息的机器人QQ
				MessageType: messageType,	//消息类型
				GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
				UserID:      userID,		//消息来源ID
				NoticID:     noticID,		//触发对象、被动
				Message:     message,		//消息内容
				RawMessage:  rawMessage,	//原始信息
				Name: 		 nick,
			}
			OnIviteRobotToGroup(&xe)
		}()
	case 215:
		//user_id邀请人
		//notic_id被邀请人
		go func() {
			nick := core.GetNick(selfID,noticID)
			xe := XEvent {
				SelfID:      selfID,		//接收消息的机器人QQ
				MessageType: messageType,	//消息类型
				GroupID:     groupID,		//消息来源ID 可能需要和下面的配合
				UserID:      userID,		//消息来源ID
				NoticID:     noticID,		//触发对象、被动
				Message:     message,		//消息内容
				RawMessage:  rawMessage,	//原始信息
				Name: 		 nick,
			}
			OnGroupMemberIntive(&xe)
		}()
	case 219:
		//这个是邀请人进群后直接进群的好像
		//notic_id是被邀请人的QQ
		//user_id是邀请人的QQ
	default:
		//
	}
	return 0
}

func XQDestroyPlugin() int64 {
	return 0
}

func XQSetUp() int64 {
	return 0
}


