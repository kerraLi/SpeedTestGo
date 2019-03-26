package service

import (
	"encoding/json"
	"fmt"
	"main/util"
	"strconv"
	"time"
)

var check *Check

type Check struct{}

func TimeCheck(interval int) *Check {
	if check == nil {
		check = &Check{}
		check.init(interval)
	}
	return check
}

func (c Check) Manage() {
	fmt.Println("manage...")
}

func (c Check) init(interval int) {
	t := time.NewTicker(time.Duration(interval*1000) * time.Millisecond)
	// 延迟执行stop:清理释放资源
	defer t.Stop()
	for {
		// 缓存通道
		expire := <-t.C
		timestamp := strconv.FormatInt(expire.Unix(), 10)

		context, _ := json.Marshal(map[string]interface{}{
			"timestamp": timestamp,
		})
		result := map[string]string{
			"action":      "monitor_time",
			"server_name": util.GetConfig("server.name"),
			"server_ip":   util.GetConfig("server.ip"),
			"result":      string(context),
		}
		msg, _ := json.Marshal(result)
		util.HandlePush("monitor", msg)
	}
}
