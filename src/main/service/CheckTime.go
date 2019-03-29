package service

import (
	"encoding/json"
	"main/util"
	"strconv"
	"time"
)

var check *Check

func CheckTime(interval int) *Check {
	if check == nil {
		check = &Check{}
		check.init(interval, TimeHandle)
	}
	return check
}

func TimeHandle() {
	action := "monitor_time"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	context, _ := json.Marshal(map[string]interface{}{
		"timestamp": timestamp,
	})
	result := map[string]string{
		"action":      action,
		"server_name": util.GetConfig("server.name"),
		"server_ip":   util.GetConfig("server.ip"),
		"result":      string(context),
	}
	queueChan <- result
}
