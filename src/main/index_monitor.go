package main

import (
	"main/service"
	"main/util"
	"strconv"
)

func main() {
	interval, _ := strconv.Atoi(util.GetConfig("monitor.interval"))
	service.TimeCheck(interval)
}
