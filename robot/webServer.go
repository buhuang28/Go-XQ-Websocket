package robot

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func WebStart()  {
	router := gin.Default()
	router.POST("/test", apiTest)
	router.Run(":8885")
}

//HttpApi接口
func apiTest(context *gin.Context)  {
	body, err := ioutil.ReadAll(context.Request.Body)
	if body == nil || err != nil {
		context.JSON(200, gin.H{
			"code":    -1,
			"message": "数据丢失",
		})
		return
	}
	var data ApiData
	err = json.Unmarshal(body, &data)
	if err != nil {
		context.JSON(200, Result{Code : -1,Message: "信息序列化失败"})
		return
	}
	switch data.ApiType {
	case 1:
		//发送私聊API
		message := SendPrivateMessage(data.RobotQQ, data.MessageType,data.GroupID, data.UserID, data.Message)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: message})
	case 2:
		//发送群聊API
		message := SendGroupMessage(data.RobotQQ, data.MessageType,data.MsgID, data.GroupID,data.UserID, data.Message)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: message})
	case 3:
		//踢人API
		KickMember(data.RobotQQ,data.GroupID,data.UserID,data.RejectMsg)
	case 4:
		//
		message := RecallMessage(data.RobotQQ, data.GroupID,data.MessageNum,data.MsgID)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: message})
	case 5:
		BanMember(data.RobotQQ,data.GroupID,data.UserID,data.Time)
	case 6:
		RejectMember(data.RobotQQ,data.SubType,data.UserID,data.GroupID,data.Approve,data.RawMessage,data.RejectMsg)
	case 7:
		groups := GetGroups(data.RobotQQ)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: groups})
	case 8:
		members := GetGroupMembers(data.RobotQQ, data.GroupID)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: members})
	case 9:
		qs := GetOnlineQQs()
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: qs})
	case 10:
		members := GetAllGroupMembers(data.RobotQQ)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: members})
	case 11:
		members_2 := GetGroupMembers_2(data.RobotQQ, data.GroupID)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: members_2})
	case 12:
		members_2 := GetAllGroupMembers_2(data.RobotQQ)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: members_2})
	case 13:
		members_2 := GetPointGroupMembers_2(data.RobotQQ, data.GroupsID)
		context.JSON(200, Result{Code : -1,Message: "sucess",Data: members_2})
	}
}

type Result struct {
	Code int `json:"code"`
	Data string `json:"data"`
	Message string `json:"message"`
}