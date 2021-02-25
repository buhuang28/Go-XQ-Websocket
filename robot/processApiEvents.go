package robot

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"xq-go-sdk/core"
)

//发送私聊
//机器人QQ
//好友/群临时会话
//接受的QQ
//消息内容
func SendPrivateMessage(robotQQ,messageType,groupID,userID int64,msg string) string {
	//if len(pics) > 0 {
	//	for _,pic := range pics {
	//		if pic != "" {
	//			var fileByte []byte
	//			var err error
	//			if strings.Contains(pic,"http"){
	//				fileByte = DownPic(pic)
	//			}else if strings.Contains(pic,`:\`) {
	//				fileByte, err = ioutil.ReadFile(pic)
	//			}
	//			if err == nil {
	//				pic = core.UpLoadPic(robotQQ, messageType, userID, fileByte)
	//			}
	//		}
	//		msg += pic
	//	}
	//}
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			logger.Println(err)
		}
	}()

	v2 := core.SendMsgEX_V2(robotQQ, messageType, groupID, userID, msg, 0, false, "")
	return v2
}

//发送群聊
//机器人QQ
//群/讨论组
//群ID
//消息内容
func SendGroupMessage(robotQQ,messageType,msgID,groupID,userID int64,msg string) string {
	//if len(pics) > 0 {
	//	for _,pic := range pics {
	//		if pic != "" {
	//			var fileByte []byte
	//			var err error
	//			if strings.Contains(pic,"http"){
	//				fileByte = DownPic(pic)
	//			}else if strings.Contains(pic,`:\`) {
	//				fileByte, err = ioutil.ReadFile(pic)
	//			}
	//			if err == nil {
	//				pic = core.UpLoadPic(robotQQ, messageType, userID, fileByte)
	//			}
	//		}
	//		msg += pic
	//	}
	//}
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			logger.Println(err)
		}
	}()

	mid := msgID
	v2 := core.SendMsgEX_V2(robotQQ, messageType, groupID, userID, msg, 0, false, "")
	var msgPack MsgPack
	err := json.Unmarshal([]byte(v2), &msgPack)
	if err == nil {
		msgNumToGroup := make(map[int64]int64)
		msgLock.Lock()
		if msgDiyIDToNum[mid] != nil && len(msgDiyIDToNum[mid]) > 0 {
			msgNumToGroup = msgDiyIDToNum[mid]
		}
		msgNumToGroup[msgPack.Msgno] = groupID
		msgDiyIDToNum[mid] = msgNumToGroup
		msgLock.Unlock()
	}
	return v2
}

//踢人
//机器人QQ
//踢出群号
//踢出QQ
//是否再次接受申请，真为不再接收，假为接收
//----感觉无需二次封装
func KickMember(robotQQ,groupID,userID int64,kickMsg string) {
	if kickMsg != "" {
		SendPrivateMessage(robotQQ,4,groupID,userID,kickMsg)
	}
	core.KickGroupMBR(robotQQ,groupID,userID,false)
}

//撤回消息
//机器人QQ
//群ID
//消息序号
//消息ID
func RecallMessage(robotQQ,groupID,msgNum,msgID int64) string {

	if msgNum != 0 && groupID != 0 && msgID != 0{
		if logger != nil {
			logger.Println("使用群号撤回")
		}
		core.WithdrawMsg(robotQQ, groupID, msgNum, msgID)
		return ""
	}

	messageNums := msgDiyIDToNum[msgID]
	if messageNums == nil || len( messageNums) == 0 {
		if logger != nil {
			logger.Println("没有messageNums")
		}
		return ""
	}
	//k是消息num，v是群号
	for k,v := range messageNums {
		messageID := msgNumToID[k]

		recallMsg := ""
		if groupID != 0 {
			recallMsg = core.WithdrawMsg(robotQQ, groupID, k, messageID)
			logger.Println("非正常撤回撤回群:",groupID,"序号:",k,"消息ID为",messageID)
		}else {
			recallMsg = core.WithdrawMsg(robotQQ, v, k, messageID)
			logger.Println("正常撤回撤回群:",v,"序号:",k,"消息ID为",messageID)
		}
		logger.Println("撤回返回数据:",recallMsg)
	}

	return ""
}

//禁言
func BanMember(robotQQ,groupID,memberID,time int64) {
	core.ShutUP(robotQQ,groupID,memberID,time)
}

//拒绝加群
func RejectMember(robotQQ,subType,userID,groupID,approve int64,rawMessage,rejectMsg string)  {
	core.HandleGroupEvent(robotQQ,subType,userID,groupID,rawMessage,approve,rejectMsg)
}

//获取群列表
func GetGroups(robotQQ int64) string {
	groups := core.GetGroupList(robotQQ)
	groups = strings.TrimSpace(groups)
	if groups != "" {
		var getGroupInfos GetGroupInfos
		err := json.Unmarshal([]byte(groups), &getGroupInfos)
		if err != nil {
			if logger != nil {
				logger.Println("groups反序列化失败:",groups)
			}
			return groups
		}
		if getGroupInfos.Ec != 0 || getGroupInfos.Errcode != 0 {
			if logger != nil {
				logger.Println("getGroupInfos信息有问题:",groups)
			}
			return groups
		}
		for _,v := range getGroupInfos.Manage {
			getGroupInfos.Join = append(getGroupInfos.Join,v)
		}

		var apiCallBack ApiCallBack
		apiCallBack.Type = GET_GROUPS
		apiCallBack.Robot = robotQQ
		apiCallBack.GroupInfos = getGroupInfos.Join
		marshal, err := json.Marshal(apiCallBack)
		if err == nil {
			WsCon.Write(marshal)
		}else {
			if logger != nil {
				logger.Println("获取QQ群信息反序列化失败:",groups)
			}
		}
	}
	return groups
}

//获取钦点群群员列表
func GetGroupMembers(robotQQ,groupID int64) string {
	memberList := core.GetGroupMemberList_C(robotQQ, groupID)
	memberList = strings.TrimSpace(memberList)
	membersRunes := []rune(memberList)
	memberRunelen := len(membersRunes)
	for string(membersRunes[memberRunelen -1]) != "}"{
		memberRunelen--
	}
	memberList = string(membersRunes[:memberRunelen])
	if memberList != "" {
		var data GetMemberData
		err := json.Unmarshal([]byte(memberList), &data)
		if err != nil {
			logger.Println("memberList反序列化失败:",memberList)
			return memberList
		}
		Mmebers := make(map[string][]Member)
		Mmebers[strconv.FormatInt(groupID,10)] = data.List
		var apiCallBack ApiCallBack
		apiCallBack.Type = GET_GROUP_MEMBERS
		apiCallBack.Robot = robotQQ
		apiCallBack.Members = Mmebers
		marshal, err := json.Marshal(apiCallBack)
		if err == nil {
			WsCon.Write(marshal)
		}else {
			logger.Println("获取QQ群成员信息反序列化失败:",memberList)
		}
	}
	return memberList
}

//获取钦点的群成员列表2
func GetGroupMembers_2(robotQQ,groupID int64) string {
	memberList := core.GetGroupMemberList_B(robotQQ, groupID)
	memberList = strings.TrimSpace(memberList)
	membersRunes := []rune(memberList)
	memberRunelen := len(membersRunes)
	for string(membersRunes[memberRunelen -1]) != "}"{
		memberRunelen--
	}
	memberList = string(membersRunes[:memberRunelen])
	if memberList != "" {
		var data GetMemberData
		err := json.Unmarshal([]byte(memberList), &data)
		if err != nil {
			logger.Println("memberList反序列化失败:",memberList)
			return memberList
		}
		Mmebers := make(map[string][]Member)
		Mmebers[strconv.FormatInt(groupID,10)] = data.List
		var apiCallBack ApiCallBack
		apiCallBack.Type = GET_GROUP_MEMBERS_2
		apiCallBack.Robot = robotQQ
		apiCallBack.Members = Mmebers
		marshal, err := json.Marshal(apiCallBack)
		if err == nil {
			WsCon.Write(marshal)
		}else {
			logger.Println("获取QQ群成员信息反序列化失败:",memberList)
		}
	}
	return memberList
}

//获取在线QQ列表
func GetOnlineQQs() string {
	list := core.GetOnLineList()
	list = strings.TrimSpace(list)
	var QQs []string
	var apiCallBack ApiCallBack
	apiCallBack.Type = GET_ONLINE_TYPE
	if list != "" {
		if strings.Contains(list,"\r\n") {
			split := strings.Split(list, "\r\n")
			for _,v := range split {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				_, err := strconv.ParseInt(v, 10, 64)
				for err != nil {
					runes := []rune(v)
					if len(runes) < 5 {
						return ""
					}
					v = string(runes[:len(runes) - 1])
					_, err = strconv.ParseInt(v, 10, 64)
				}
				QQs = append(QQs,v)
			}
		}else {
			_, err := strconv.ParseInt(list, 10, 64)
			for err != nil {
				runes := []rune(list)
				if len(runes) < 5 {
					return ""
				}
				list = string(runes[:len(runes) - 1])
				if list == "" {
					return ""
				}
				_, err = strconv.ParseInt(list, 10, 64)
			}
			QQs = append(QQs,list)
		}
	}else {
		return ""
	}

	list2 := core.GetOnLineList()
	list2 = strings.TrimSpace(list2)
	var QQs2 []string
	if list2 != "" {
		if strings.Contains(list2,"\r\n") {
			split := strings.Split(list2, "\r\n")
			for _,v := range split {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				_, err := strconv.ParseInt(v, 10, 64)
				for err != nil {
					runes := []rune(v)
					if len(runes) < 5 {
						return ""
					}
					v = string(runes[:len(runes) - 1])
					_, err = strconv.ParseInt(v, 10, 64)
				}
				QQs2 = append(QQs2,v)
			}
		}else {
			_, err := strconv.ParseInt(list2, 10, 64)
			for err != nil {
				runes := []rune(list2)
				if len(runes) < 5 {
					return ""
				}
				list2 = string(runes[:len(runes) - 1])
				if list2 == "" {
					return ""
				}
				_, err = strconv.ParseInt(list2, 10, 64)
			}
			QQs2 = append(QQs2,list2)
		}
	}else {
		return ""
	}

	if len(QQs) != len(QQs2) {
		logger.Println("QQs:",QQs,"\r\n","QQs2",QQs2)
		return ""
	}

	var trueQQs []string

	for k,v := range QQs {
		if QQs[k] == QQs2[k] {
			trueQQs = append(trueQQs,v)
			continue
		}

		var trueRune []rune

		Q1Runes := []rune(QQs[k])
		Q2Runes := []rune(QQs2[k])
		if len(Q1Runes) != len(Q2Runes) {
			logger.Println("QQs:",QQs,"\r\n","QQs2",QQs2)
			if len(Q1Runes) > len(Q2Runes) {
				trueRune = Q2Runes
			}else {
				trueRune = Q1Runes
			}
			trueQQs = append(trueQQs,string(trueRune))
			continue
		}

		for k,v := range Q1Runes {
			if Q1Runes[k] == Q2Runes[k] {
				trueRune = append(trueRune,v)
			}else {
				logger.Println("QQs:",QQs,"\r\n","QQs2",QQs2)
				break
			}
		}
		trueQQs = append(trueQQs,string(trueRune))
	}

	apiCallBack.OnLines = trueQQs

	marshal, err := json.Marshal(apiCallBack)
	if err == nil {
		if WsCon != nil {
			WsCon.Write(marshal)
		}
	}else {
		if logger != nil {
			logger.Println("获取QQ列表json反序列化失败:",list)
		}
	}
	qs := ""
	for _,v := range trueQQs {
		qs += v+","
	}
	qs = strings.Trim(qs,",")
	return qs
}

//获取全部群和成员
func GetAllGroupMembers(robotQQ int64) string {

	var data ApiCallBack

	GroupMemberMap := make(map[string][]Member)

	groupList := core.GetGroupList(robotQQ)
	groupList = strings.TrimSpace(groupList)
	if groupList == "" {
		return groupList
	}

	parse := gjson.Parse(groupList)

	Ec := parse.Get("ec").Int()
	Er := parse.Get("errcode").Int()
	if Ec != 0 || Er != 0 {
		logger.Println("getGroupInfos信息有问题:",groupList)
		return groupList
	}

	joinArray := parse.Get("join").Array()
	manageArray := parse.Get("manage").Array()


	var managerGroup []int64
	for _,v := range manageArray {
		joinArray = append(joinArray,v)
		managerGroup = append(managerGroup,v.Get("gc").Int())
	}

	for _,v := range joinArray {
		membersString := core.GetGroupMemberList_C(robotQQ, v.Get("gc").Int())
		membersString = strings.TrimSpace(membersString)

		membersJSONValue := gjson.Parse(membersString)
		membersJSONArray := membersJSONValue.Get("list").Array()
		var memberList []Member
		for _,v2 := range membersJSONArray {
			var member Member
			member.QQ = strconv.FormatInt(v2.Get("QQ").Int(),10)
			member.Lv = v2.Get("lv").Int()
			member.Val = v2.Get("val").Int()
			memberList = append(memberList,member)
		}
		GroupMemberMap[strconv.FormatInt(v.Get("gc").Int(),10)] = memberList
	}

	if len(GroupMemberMap) == 0 {
		logger.Println(robotQQ,"获取不到QQ群和成员信息")
		return "获取不到QQ群和成员信息"
	}

	data.Type = GET_ALL_MEMBERS
	data.Robot = robotQQ
	data.Members = GroupMemberMap
	data.ManagerGroup = managerGroup
	marshal, err := json.Marshal(data)
	if err != nil {
		logger.Println("data序列化失败")
		return "data序列化失败"
	}
	if WsCon != nil {
		WsCon.Write(marshal)
	}
	return string(marshal)
}

//获取全部群和成员2
func GetAllGroupMembers_2(robotQQ int64) string {
	var data ApiCallBack

	GroupMemberMap := make(map[string]map[string]Nick)

	groups := core.GetGroupList(robotQQ)
	groups = strings.TrimSpace(groups)
	if groups == "" {
		return groups
	}
	runes := []rune(groups)
	end := len(runes)
	for string(runes[end -1]) != "}"{
		end--
	}
	groups = string(runes[:end])
	var getGroupInfos GetGroupInfos
	err := json.Unmarshal([]byte(groups), &getGroupInfos)
	if err != nil {
		logger.Println("GetAllGroupMembers_2反序列化失败:",groups)
		return groups
	}
	if getGroupInfos.Ec != 0 || getGroupInfos.Errcode != 0 {
		logger.Println("GetAllGroupMembers_2信息有问题:",groups)
		return groups
	}
	if getGroupInfos.Join == nil || len(getGroupInfos.Join) == 0 {
		logger.Println("GetAllGroupMembers_2里无法获取到群信息:",groups)
		return groups
	}

	var manageGroups []int64

	for _,v := range getGroupInfos.Manage {
		getGroupInfos.Join = append(getGroupInfos.Join,v)
		manageGroups = append(manageGroups,v.Gc)
	}

	for _,v := range getGroupInfos.Join {
		members := core.GetGroupMemberList_B(robotQQ, v.Gc)
		members = strings.TrimSpace(members)
		membersRunes := []rune(members)
		memberRunelen := len(membersRunes)
		for string(membersRunes[memberRunelen -1]) != "}"{
			memberRunelen--
		}
		members = string(membersRunes[:memberRunelen])
		if members == "" {
			logger.Println(v.Gc,"获取不到QQ成员信息")
			continue
		}
		var data Member2
		err := json.Unmarshal([]byte(members), &data)
		if err != nil {
			logger.Println(v.Gc,"的members反序列化失败:",members)
			continue
		}
		GroupMemberMap[strconv.FormatInt(v.Gc,10)] = data.Members
	}

	if len(GroupMemberMap) == 0 {
		logger.Println(robotQQ,"获取不到QQ群和成员信息")
		return groups
	}
	data.Type = GET_ALL_MEMBERS_2
	data.Robot = robotQQ
	data.Members2 = GroupMemberMap
	data.ManagerGroup = manageGroups
	marshal, err := json.Marshal(data)
	if err != nil {
		logger.Println("data序列化失败")
		return groups
	}
	if WsCon != nil {
		WsCon.Write(marshal)
	}
	return string(marshal)
}

func GetPointGroupMembers_2(robotQQ int64,pointGroup []string) string {
	var data ApiCallBack
	groups := core.GetGroupList(robotQQ)
	groups = strings.TrimSpace(groups)
	if groups == "" {
		return groups
	}
	runes := []rune(groups)
	end := len(runes)
	for string(runes[end -1]) != "}"{
		end--
	}
	groups = string(runes[:end])
	var getGroupInfos GetGroupInfos
	err := json.Unmarshal([]byte(groups), &getGroupInfos)
	if err != nil {
		logger.Println("groups反序列化失败:",groups)
		return groups
	}
	if getGroupInfos.Ec != 0 || getGroupInfos.Errcode != 0 {
		logger.Println("getGroupInfos信息有问题:",groups)
		return groups
	}
	if getGroupInfos.Join == nil || len(getGroupInfos.Join) == 0 {
		logger.Println("getGroupInfos里无法获取到群信息:",groups)
		return groups
	}
	for _,v := range getGroupInfos.Manage {
		getGroupInfos.Join = append(getGroupInfos.Join,v)
	}

	//群--QQ号--昵称
	GroupMemberMap := make(map[string]map[string]Nick)
	for _,v := range getGroupInfos.Join {
		groupString := strconv.FormatInt(v.Gc, 10)
		logger.Println("groupString:",groupString)
		if !Contain(groupString,pointGroup) {
			continue
		}
		members := core.GetGroupMemberList_B(robotQQ, v.Gc)
		members = strings.TrimSpace(members)
		membersRunes := []rune(members)
		memberRunelen := len(membersRunes)
		for string(membersRunes[memberRunelen -1]) != "}"{
			memberRunelen--
		}
		members = string(membersRunes[:memberRunelen])
		if members == "" {
			logger.Println(v.Gc,"获取不到QQ成员信息")
			continue
		}
		var data Member2
		err := json.Unmarshal([]byte(members), &data)
		if err != nil {
			logger.Println(v.Gc,"的members反序列化失败:",members)
			continue
		}
		GroupMemberMap[strconv.FormatInt(v.Gc,10)] = data.Members
	}

	if len(GroupMemberMap) == 0 {
		logger.Println(robotQQ,"获取不到QQ群和成员信息")
		return groups
	}
	data.Type = GET_PONIT_MEMBERS_2
	data.Robot = robotQQ
	data.Members2 = GroupMemberMap
	marshal, err := json.Marshal(data)
	if err != nil {
		logger.Println("data序列化失败")
		return groups
	}
	if WsCon != nil {
		WsCon.Write(marshal)
	}
	return string(marshal)
}
