package controller

import (
	"fmt"
	"main/service"
	"net/http"
)

func SpeedTest(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		url := req.FormValue("url")
		fmt.Println(url)
		service.SpeedUrl(url)
	}
}
