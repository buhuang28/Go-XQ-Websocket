package robot

import "encoding/json"

func processWebsocketMsg(websocketMsg []byte)  {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("GetOnlineQQs异常:",err)
			}else {
				writeFile("exc.txt","onStart异常")
			}
		}
	}()
	var data ApiData
	var err error
	err = json.Unmarshal(websocketMsg, &data)
	if err != nil {
		if logger != nil {
			logger.Println("反序列化失败",string(websocketMsg))
		}
		return
	}else {
		//if logger != nil {
		//	logger.Println(string(websocketMsg))
		//}
	}
	switch data.ApiType {
	case 1:
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
			SendPrivateMessage(data.RobotQQ, data.MessageType,data.GroupID, data.UserID, data.Message)
		}()
	case 2:
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
			SendGroupMessage(data.RobotQQ, data.MessageType,data.MsgID, data.GroupID,data.UserID, data.Message)
		}()
	case 3:
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
			KickMember(data.RobotQQ,data.GroupID,data.UserID,data.RejectMsg)
		}()
	case 4:
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
			RecallMessage(data.RobotQQ, data.GroupID,data.MessageNum, data.MsgID)
		}()
	case 5:
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
			BanMember(data.RobotQQ,data.GroupID,data.UserID,data.Time)
		}()
	case 6:
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
			RejectMember(data.RobotQQ,data.SubType,data.UserID,data.GroupID,data.Approve,data.RawMessage,data.RejectMsg)
		}()
	case 7:
		GetGroups(data.RobotQQ)
	case 8:
		GetGroupMembers(data.RobotQQ, data.GroupID)
	case 9:
		GetOnlineQQs()
	case 10:
		GetAllGroupMembers(data.RobotQQ)
	case 11:
		GetGroupMembers_2(data.RobotQQ, data.GroupID)
	case 12:
		GetAllGroupMembers_2(data.RobotQQ)
	case 13:
		GetPointGroupMembers_2(data.RobotQQ, data.GroupsID)
	}
}
