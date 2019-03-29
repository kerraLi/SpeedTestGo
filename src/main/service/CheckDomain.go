package service

import (
	"crypto/tls"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/oschwald/geoip2-golang"
	"io/ioutil"
	"main/util"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SpeedInfo struct {
	Ip                string `json:"ip"`
	Url               string `json:"url"`
	Address           string `json:"address"`
	Status            string `json:"status"`
	Msg               string `json:"msg"`
	RedirectCount     int    `json:"redirect_count"`
	HttpCode          int    `json:"http_code"`
	DnsResolveTime    string `json:"dns_resolve_time"`
	HttpConnTime      string `json:"http_conn_time"`
	HttpPreTrans      string `json:"http_pre_trans"`
	HttpStartTrans    string `json:"http_start_trans"`
	HttpTotalTime     string `json:"http_total_time"`
	HttpSizeDownload  string `json:"http_size_download"`
	HttpSpeedDownload string `json:"http_speed_download"`
}

var checkDomain *Check
var redisClient *redis.Client

func CheckDomain(interval int) *Check {
	if check == nil {
		check = &Check{}
		// monitor
		go func() {
			check.monitor()
		}()
		// begin
		check.init(interval, DomainHandle)
	}
	return checkDomain
}

func DomainHandle() {
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
		util.FailOnErrorNoExit(err, "redis 连接获取数据异常")
		return
	}
	val, _ = url.QueryUnescape(val)
	err = json.Unmarshal([] byte(val[1:len(val)-1]), &domains)
	if err != nil {
		util.FailOnErrorNoExit(err, "domains 获取异常")
		return
	}
	// speed
	ch := make(chan int, 500)
	for _, theDomain := range domains {
		ch <- 1
		// 协程
		go func(domain string) {
			speed := SpeedUrl(domain)
			if speed.Status != "success" {
				pushAlert(speed)
			}
			<-ch
		}(theDomain)
	}
	close(ch)
}

// 测速 return SpeedInfo
func SpeedUrl(domain string) SpeedInfo {

	var speed SpeedInfo
	var start time.Time
	var ip net.IPAddr
	// parse
	testUrl, err := parseURL(domain)
	if err != nil {
		speed.Status = "error"
		speed.Msg = err.Error()
		return speed
	}
	speed.Url = testUrl.String()
	speed.RedirectCount = -1

	// req
	req, _ := http.NewRequest("GET", testUrl.String(), nil)
	trace := &httptrace.ClientTrace{
		// dns解析时间
		DNSStart: func(dsi httptrace.DNSStartInfo) {},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			if ddi.Err == nil {
				ip = ddi.Addrs[0]
			}
			speed.DnsResolveTime = time.Since(start).String()
		},

		// tls握手时间
		TLSHandshakeStart: func() {},
		TLSHandshakeDone:  func(cs tls.ConnectionState, err error) {},

		// 建立连接时间
		ConnectStart: func(network, addr string) {},
		ConnectDone: func(network, addr string, err error) {
			speed.RedirectCount = speed.RedirectCount + 1
			speed.HttpConnTime = time.Since(start).String()
		},

		// 获取到连接
		GotConn: func(info httptrace.GotConnInfo) {
			speed.HttpPreTrans = time.Since(start).String()
		},

		// 获取首字节时间(server_processing)
		GotFirstResponseByte: func() {
			speed.HttpStartTrans = time.Since(start).String()
		},
	}

	// req
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	// transport
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
		// 是否追踪301跳转：允许
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	// always refuse to follow redirects, visit does that
		//	// manually if required.
		//	return http.ErrUseLastResponse
		//},
	}
	req.Header.Set("Cache-Control", "no-cache")
	// resp
	resp, err := client.Do(req)
	if err != nil {
		speed.Status = "error"
		speed.HttpTotalTime = time.Since(start).String()
		speed.Msg = err.Error()
		if ip.IP != nil {
			speed.Ip = ip.IP.String()
			speed.Address = getIpLocation(ip.IP)
		}
		return speed
	}
	defer resp.Body.Close()
	// read
	body, _ := ioutil.ReadAll(resp.Body)
	// data
	speed.Status = "success"
	speed.HttpCode = resp.StatusCode
	speed.HttpTotalTime = time.Since(start).String()
	speed.HttpSizeDownload = strconv.Itoa(len(body))
	speed.HttpSpeedDownload = ""
	speed.Ip = ip.IP.String()
	speed.Address = getIpLocation(ip.IP)
	if speed.HttpCode != 200 {
		speed.Status = "error"
	}
	return speed
}

// 处理scheme
func parseURL(uri string) (*url.URL, error) {
	if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
		uri = "//" + uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return u, err
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u, err
}

// 获取ip地址
func getIpLocation(ip net.IP) string {
	if ip == nil {
		return ""
	}
	db, err := geoip2.Open("src/lib/GeoLite2-City.mmdb")
	if err != nil {
		util.FailOnErrorNoExit(err, "ip 地址库异常")
		return ""
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	record, err := db.City(ip)
	if err != nil {
		util.FailOnErrorNoExit(err, "ip 地址搜索异常")
		return ""
	}
	//fmt.Println(record.Continent, record.Country, record.Subdivisions, record.City)

	continent := record.Continent.Names["zh-CN"]
	country := record.Country.Names["zh-CN"]
	if len(record.Subdivisions) == 0 {
		return continent + country
	}
	subdivisions := record.Subdivisions[0].Names["zh-CN"]
	city := record.City.Names["zh-CN"]
	return continent + country + subdivisions + city
}

// 报警推送
func pushAlert(info SpeedInfo) {
	action := "monitor_domain"
	context, _ := json.Marshal(info)
	result := map[string]string{
		"action":      action,
		"server_name": util.GetConfig("server.name"),
		"server_ip":   util.GetConfig("server.ip"),
		"result":      string(context),
	}
	queueChan <- result
}
