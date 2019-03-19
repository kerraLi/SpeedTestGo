package main

import (
	"./controller"
	"./service"
	"./util"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 初始化
	args := Args{}
	args.init()

	// 常见的调用方式
	flag.Parse()
	if args.usage {
		flag.Usage()
	}
	args.Handler()

}

// Param存储接收到的标识
type Args struct {
	usage         bool
	configHandler string
	timeHandler   string
	timeInterval  int
}

func (args *Args) len() int {
	count := 0

	return count
}

// init 根据传入的标识初始化配置
func (args *Args) init() {
	flag.BoolVar(&args.usage, "h", false, "help information")
	flag.StringVar(&args.configHandler, "c", "", "ConfigManage: start/stop/restart set config manage service `option` [port:8888]")
	flag.StringVar(&args.timeHandler, "t", "", "TimeMonitor: start/stop/restart time monitor service `option`")
	flag.IntVar(&args.timeInterval, "s", 60, "TimeMonitor: the time interval of time handler[unit:`second`]")
	//name
	name := util.GetConfig("main.name")
	version := util.GetConfig("main.version")
	//usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%s version: %s
Usage: %s [-h] [-c option] [-t option][-s second]

Options:
`, name, version, name)
		flag.PrintDefaults()
	}
}

// handler
func (args *Args) Handler() {
	// config manage
	switch args.configHandler {
	case "start":
		//绑定路由 如果访问 /upload 调用 Handler 方法
		http.HandleFunc("/upload", controller.ConfigManage)
		//使用 tcp 协议监听8888
		http.ListenAndServe(":8888", nil)
	case "stop":
	case "restart":
	}

	// time monitor
	switch args.timeHandler {
	case "start":
		service.TimeCheck(args.timeInterval)
	case "stop":
	case "restart":
	}

	flag.Usage()
}
