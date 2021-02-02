package robot

import (
	"log"
	"sync"
)

var (
	//通过自己定义的ID存储群消息返回的  Num----对应的群号
	msgDiyIDToNum = make(map[int64]map[int64]int64)
	//通过消息Num得到ID
	msgNumToID = make(map[int64]int64)

	//全局日志
	logger *log.Logger

	//消息Map锁
	msgLock sync.Mutex
)

//websocket、httpApi调用包(被调用)
type ApiData struct {
	ApiType int64 `json:"api_type"`
	MsgID int64 `json:"msg_id"`
	RobotQQ int64 `json:"robot_qq"`
	Message string `json:"message"`
	MessageType int64 `json:"message_type"`
	UserID int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	MemberID int64 `json:"member_id"`
	MessageNum int64 `json:"message_num"`
	MessageID int64 `json:"message_id"`
	Time int64 `json:"time"`
	SubType int64 `json:"sub_type"`
	RawMessage string `json:"raw_message"`
	RejectMsg string `json:"reject_msg"`
	Approve int64 `json:"approve"`
	//Members map[string][]Member `json:"members"`
	GroupsID []string `json:"groups_id"`
}

//发送消息的返回数据包
type MsgPack struct {
	Packcount int64    `json:"packcount"`
	Msgid     int64    `json:"msgid"`
	Msgno     int64    `json:"msgno"`
	Msgtime   int64    `json:"msgtime"`
	Sendok    bool     `json:"sendok"`
	Txtinf    string   `json:"txtinf"`
}

//主动返回给websocket-Server
type ApiCallBack struct {
	Type int64 `json:"type"`
	Robot int64 `json:"robot"`
	OnLines []string `json:"on_lines"`
	GroupInfos []GroupInfo `json:"group_infos"`
	Members map[string][]Member `json:"members"`
	Members2 map[string]map[string]Nick `json:"members_2"`
	ManagerGroup []int64 `json:"manager_group"`
}

type GetMemberData struct {
	List []Member `json:"list"`
}

type Member struct {
	QQ  string `json:"QQ"`
	Lv  int64    `json:"lv"`
	Val int64    `json:"val"`
}

type Member2 struct {
	Ec       int64    `json:"ec"`
	Errcode  int64    `json:"errcode"`
	Em       string `json:"em"`
	C        string `json:"c"`
	ExtNum   int64    `json:"ext_num"`
	Level    int64    `json:"level"`
	MemNum   int64    `json:"mem_num"`
	MaxNum   int64    `json:"max_num"`
	MaxAdmin int64    `json:"max_admin"`
	Owner    int64    `json:"owner"`
	Members map[string]Nick `json:"members"`
	Adm []int64 	`json:"adm"`
}

type Nick struct {
	Nk string `json:"nk"`
}
//
//type Nick struct {
//	Lst int    `json:"lst"`
//	Jt  int    `json:"jt"`
//	Rm  int    `json:"rm"`
//	Lad int    `json:"lad"`
//	Lp  int    `json:"lp"`
//	Ll  int    `json:"ll"`
//	Nk  string `json:"nk"`
//}

type GetGroupInfos struct {
	Ec      int64    `json:"ec"`
	Errcode int64    `json:"errcode"`
	Em      string `json:"em"`
	Join    []GroupInfo `json:"join"`
	Manage  []GroupInfo `json:"manage"`
}

type GroupInfo struct {
	Gc    int64    `json:"gc"`
	Gn    string `json:"gn"`
	Owner int64    `json:"owner"`
}

const (
	GET_PONIT_MEMBERS_2 int64 = 13
	GET_ALL_MEMBERS_2 int64 = 12
	GET_GROUP_MEMBERS_2 int64 = 11
	GET_ALL_MEMBERS int64 = 10
	GET_ONLINE_TYPE int64 = 9
	GET_GROUP_MEMBERS int64 = 8
	GET_GROUPS int64 = 7
	REJECT_REQUEST int64 = 6
	BAN_MEMBER int64 = 5
	RECALL_MESSAGE int64 = 4
	KICK_MEMBER int64 = 3
	SEND_GROUP_MESSAGE int64 = 2
	SEND_PRIVATE_MESSAGE int64 = 1
)