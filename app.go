package main

import (
	"encoding/json"
	"xq-go-sdk/robot"
)

// 插件信息
type AppInfo struct {
	Name   string `json:"name"`   // 插件名字
	Pver   string `json:"pver"`   // 插件版本
	Sver   int    `json:"sver"`   // 框架版本
	Author string `json:"author"` // 作者名字
	Desc   string `json:"desc"`   // 插件说明
}

func newAppInfo() *AppInfo {
	return &AppInfo {
		Name:   "Websocket-Client",
		Pver:   "1.3.36",
		Sver:   3,
		Author: "不慌",
		Desc:   "websocket-Client",
	}
}

func init() {
	data, _ := json.Marshal(newAppInfo())
	robot.AppInfoJson = string(data)
}

func main() { }
