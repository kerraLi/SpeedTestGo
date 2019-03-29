package service

import (
	"fmt"
	"main/util"
	"time"
)

var queueChan = make(chan map[string]string)

type Check struct{}

func (c Check) Manage() {
	fmt.Println("manage...")
}

func (c Check) init(interval int, callback func()) {
	// ticker
	t := time.NewTicker(time.Duration(interval*1000) * time.Millisecond)
	// 延迟执行stop:清理释放资源
	defer t.Stop()
	// 第一次执行
	callback()
	for {
		// 缓存通道
		<-t.C
		// 调用传入函数
		callback()
	}
}

func (c Check) monitor() {
	for {
		msg := <-queueChan
		util.HandlePush(msg)
	}
}
