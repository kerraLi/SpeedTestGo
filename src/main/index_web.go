package main

import (
	"main/controller"
	"net/http"
)

func main() {
	//绑定路由 如果访问 /upload 调用 Handler 方法
	http.HandleFunc("/upload", controller.ConfigManage)
	//使用 tcp 协议监听8888
	http.ListenAndServe(":8888", nil)
}
