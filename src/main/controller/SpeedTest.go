package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"main/service"
	"main/util"
	"net/http"
	"net/url"
	"strconv"
)

var redisClient *redis.Client

func SpeedTest(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		action := req.FormValue("action")
		if action == "speed_test" {
			domain := req.FormValue("url")
			result := service.SpeedUrl(domain)
			context, _ := json.Marshal(result)
			fmt.Fprintln(w, string(context))
			return
		} else if action == "speed_monitor" {
			// redis client
			if redisClient == nil {
				db, _ := strconv.Atoi(util.GetConfig("redis.database"))
				redisClient = redis.NewClient(&redis.Options{
					Addr:     util.GetConfig("redis.host") + ":" + util.GetConfig("redis.port"),
					Password: util.GetConfig("redis.password"),
					DB:       db,
				})
			}
			// 获取domains
			var domains []string
			val, err := redisClient.Get("MONITOR_DOMAINS").Result()
			if err != nil {
				fmt.Fprintln(w, "redis 连接获取数据异常")
				return
			}
			val, _ = url.QueryUnescape(val)
			err = json.Unmarshal([] byte(val[1:len(val)-1]), &domains)
			if err != nil {
				fmt.Fprintln(w, "domains 获取异常")
				return
			}
			// speed
			var slice []service.SpeedInfo
			// ch 控制协程数量
			ch := make(chan int, 500)
			// fCh 控制是否所有协程跑完
			fCh := make(chan int, len(domains))
			for _, theDomain := range domains {
				ch <- 1
				// 协程
				go func(domain string) {
					result := service.SpeedUrl(domain)
					slice = append(slice, result)
					<-ch
					fCh <- 1
				}(theDomain)
			}
			for i := 0; i < len(domains); i++ {
				<-fCh
			}
			c, _ := json.Marshal(slice)
			fmt.Fprintln(w, string(c))
			return
		}
	}
}
