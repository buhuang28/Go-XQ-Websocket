package robot

import "encoding/json"

func processWebsocketMsg(websocketMsg []byte)  {
	var data ApiData
	var err error
	err = json.Unmarshal(websocketMsg, &data)
	if err != nil {
		logger.Println(string(websocketMsg),"反序列化失败")
		return
	}else {
		logger.Println(string(websocketMsg))
	}
	switch data.ApiType {
	case 1:
		SendPrivateMessage(data.RobotQQ, data.MessageType,data.GroupID, data.UserID, data.Message)
	case 2:
		SendGroupMessage(data.RobotQQ, data.MessageType,data.MsgID, data.GroupID,data.UserID, data.Message)
	case 3:
		KickMember(data.RobotQQ,data.GroupID,data.UserID,data.RejectMsg)
	case 4:
		RecallMessage(data.RobotQQ, data.GroupID,data.MessageNum, data.MsgID)
	case 5:
		BanMember(data.RobotQQ,data.GroupID,data.UserID,data.Time)
	case 6:
		RejectMember(data.RobotQQ,data.SubType,data.UserID,data.GroupID,data.Approve,data.RawMessage,data.RejectMsg)
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
