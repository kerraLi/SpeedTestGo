package main

import (
	"main/controller"
	"main/util"
	"net/http"
)

func main() {
	//绑定路由
	// api-配置上传
	http.HandleFunc("/upload", controller.ConfigManage)
	// api-速度测速
	http.HandleFunc("/speed", controller.SpeedTest)
	//使用 tcp 协议监听8888
	err := http.ListenAndServe(":8888", nil)
	util.FailOnError(err, "speed服务器启动失败")
}
