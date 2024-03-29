package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"
)

type temp struct {
	fileName string
	fileType string
	content  string
}

type JsonPostSample struct {
}

var temps = map[string]temp{
	"nginx-http":   {"http_proxy.conf", "nginx-http", "/usr/local/openresty/nginx/conf/web/"},
	"nginx-https":  {"https_proxy.conf", "nginx-https", "/usr/local/openresty/nginx/conf/web/"},
	"cert-key":     {"cert.key", "cert-key", "/usr/local/openresty/nginx/conf/cert.d/"},
	"cert-crt":     {"cert.crt", "cert-crt", "/usr/local/openresty/nginx/conf/cert.d/"},
	"rewrite-rule": {"rewrite.rule", "rewrite-rule", "/usr/local/openresty/nginx/conf/rule-config/"},
	"config-lua":   {"config.lua", "config-lua", "/usr/local/openresty/lualib/resty/upstream/"},
	"filebeat-yml": {"filebeat.yml", "filebeat-yml", "/opt/filebeat/"}}

//nginx-http：没有证书以http方式访问的配置文件，文件名规定http_proxy.conf，Type规定nginx-http，生产环境路径为/usr/local/openresty/nginx/conf/web/；
//nginx-https：有证书以https方式访问的配置文件，文件名规定https_proxy.conf，Type规定nginx-https，生产环境路径为/usr/local/openresty/nginx/conf/web/，必须配合crt和key使用；
//cert-key：用于https访问是的证书key，文件名规定为cert.key，Type规定为cert-key，生产环境路径为/usr/local/openresty/nginx/conf/cert.d/；
//cert-crt：用于https访问是的证书crt，文件名规定为cert.crt，Type规定为cert-crty，生产环境路径为/usr/local/openresty/nginx/conf/cert.d/；
//rewrite-rule：该规则是规定当使用https方式访问时，需要跳转的https域名，文件名规定为rewrite-rule，Type规定为rewrite-rule，生产环境路径为/usr/local/openresty/nginx/conf/rule-config/；
//config-lua：主要配置一些防御规则开关，主要修改防御CC规则,文件名规定为config.lua，Type规定为config-lua，生产环境路径为/usr/local/openresty/lualib/resty/upstream/；
//filebeat-yml：日志filebeat配置文件，文件名规定为filebeat.yml，Type规定为filebeat-yml，生产环境路径为/opt/filebeat/;

type resultTemp struct {
	code   string
	status string
	action string
}

var uuidMap = map[string]resultTemp{}

func ConfigManage(w http.ResponseWriter, req *http.Request) {
	// 创建uuid
	uu := uuid.Must(uuid.NewV4()).String()
	uuidMap[uu] = resultTemp{uu, "init", ""}

	//输出对应的 请求方式
	fmt.Println(req.Method)
	//判断对应的请求来源。如果为get 显示对应的页面
	if req.Method == "GET" {
		fmt.Fprintln(w, "不支持这种调用方式!")
	} else if req.Method == "POST" {
		fileType := req.FormValue("fileType")

		//解析 form 中的file 上传名字
		file, file_head, file_err := req.FormFile("fileName")

		if file_err != nil {
			js := make(map[string]interface{})
			js["code"] = uu
			js["status"] = 500
			js["type"] = false
			js["msg"] = file_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			fmt.Printf("%s\r\n", string(upl))
			return
		}

		if _, ok := temps[fileType]; !ok {
			js := make(map[string]interface{})
			js["code"] = uu
			js["status"] = 500
			js["type"] = false
			js["msg"] = "fileTpye_err"
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			fmt.Printf("%s\r\n", string(upl))
			return
		}

		if file_head.Filename != temps[fileType].fileName {
			js := make(map[string]interface{})
			js["code"] = uu
			js["status"] = 500
			js["type"] = false
			js["msg"] = "fileName_err"
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			fmt.Printf("%s\r\n", string(upl))
			return
		}

		file_save := temps[fileType].content + file_head.Filename
		//打开 已只读,文件不存在创建 方式打开  要存放的路径资源
		f, f_err := os.OpenFile(file_save, os.O_WRONLY|os.O_CREATE, 0666)
		if f_err != nil {
			fmt.Fprintf(w, "file open fail:%s", f_err)
			js := make(map[string]interface{})
			js["code"] = uu
			js["status"] = 500
			js["type"] = false
			js["msg"] = f_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			fmt.Printf("%s\r\n", string(upl))
		}
		//文件 copy
		_, copy_err := io.Copy(f, file)
		if copy_err != nil {
			js := make(map[string]interface{})
			js["code"] = uu
			js["status"] = 500
			js["type"] = false
			js["msg"] = copy_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			fmt.Printf("%s\r\n", string(upl))
		}
		//关闭对应打开的文件
		defer f.Close()
		defer file.Close()

		// todo one thread
		afterUpload(uu);

		//返回上传结果
		js := make(map[string]interface{})
		js["code"] = uu
		js["stats_code"] = 200
		js["status"] = uuidMap[uu].status
		js["action"] = uuidMap[uu].action
		upl, _ := json.Marshal(js)
		fmt.Fprintln(w, string(upl))
		fmt.Printf("%s\r\n", string(upl))

		if uuidMap[uu].status == "success" {
			time.AfterFunc(3*time.Second, func() {
				back := make(map[string]interface{})
				back["runResult"] = "1"
				back["id"] = uu
				bytesData, err := json.Marshal(back)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				reader := bytes.NewReader(bytesData)
				url := "http://yunweizdh.tel5678.com/api/configManage/updateState"
				request, err := http.NewRequest("POST", url, reader)

				if err != nil {
					fmt.Println(err.Error())
					return
				}
				request.Header.Set("Content-Type", "application/json;charset=UTF-8")
				client := http.Client{}
				resp, err := client.Do(request)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				respBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				//byte数组直接转成string，优化内存
				str := (*string)(unsafe.Pointer(&respBytes))
				fmt.Println(*str, request)
			})
		}

	} else { //如果有其他方式进行页面调用。http Status Code 500
		w.WriteHeader(500)
		fmt.Fprintln(w, "不支持这种调用方式!")

	}
}

func afterUpload(uu string) {
	uuidMap[uu] = resultTemp{uu, "doing", "check"}
	check := exec.Command("/usr/local/openresty/nginx/sbin/nginx", "-t")
	out, err := check.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	var status bool = strings.Contains(string(out), "successful")

	if status == true {
		uuidMap[uu] = resultTemp{uu, "doing", "rebot"}
		check := exec.Command("/usr/local/openresty/nginx/sbin/nginx", "-s", "reload")
		out, err := check.CombinedOutput()
		fmt.Printf(string(out))
		if err != nil {
			fmt.Println(err)
		}

		var status bool = strings.Contains(string(out), "started")

		if status == true {
			uuidMap[uu] = resultTemp{uu, "success", "rebot"}
		} else {
			uuidMap[uu] = resultTemp{uu, "failure", "rebot"}
		}

	} else {
		uuidMap[uu] = resultTemp{uu, "failure", "check"}
		return
	}

}
