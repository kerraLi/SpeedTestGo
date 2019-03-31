# WEB应用
## 启动说明
* 默认端口8888
* web仅包含linux文件管理及速度测试服务；建议使用平台：linux。
* 参照src/main/config.example.yml文件配置src/main/config.yml文件
* 打包：go build src/main/index_web.go
```
//打包成64位linux可执行文件
env GOOS=linux GOARCH=amd64 go build src/main/index_web.go
```
* 设置权限777
* 执行生成得可执行文件

### API：配置上传
* 作用
>对外提供配置文件上传至服务器接口，统一服务器配置管理、操作；统一管理繁杂的服务器配置。
* URL
>{host}:{port}/upload
* 方法
>POST
* 入参
```
[
    "fileName":"文件名称(String)",
    "fileType":"文件类型(String)",
]
// 参数说明
// nginx-http：没有证书以http方式访问的配置文件，fileName规定http_proxy.conf，Type规定nginx-http  
// nginx-https：有证书以https方式访问的配置文件，fileName规定https_proxy.conf，Type规定nginx-https  
// cert-key：用于https访问是的证书key，fileName规定为cert.key，Type规定为cert-key  
// cert-crt：用于https访问是的证书crt，fileName规定为cert.crt，Type规定为cert-crty  
// rewrite-rule：该规则是规定当使用https方式访问时，需要跳转的https域名，fileName规定为rewrite-rule，Type规定为rewrite-rule  
// config-lua：主要配置一些防御规则开关，主要修改防御CC规则,fileName规定为config.lua，Type规定为config-lua  
// filebeat-yaml：日志filebeat配置文件，fileName规定为filebeat.yaml，Type规定为filebeat-yaml 
```
* 出参
```
[
    "status":"状态 错误状态500",
    "type":"标记为false",
    "msg":"错误内容信息",
]
```

### API：域名速度检测

* 作用
>服务器作为测速节点，对外提供域名速度检测接口，快速检测域名再当前服务器区域解析结果、速度等。
* URL
>{host}:{port}/speed
* 方法
>POST
* 入参
```
[
    "action":"speed-test单域名速度检测/speed-monitor多域名速度批量检测",
    "url":"action为speed-test时，必须传入url参数，该参数为待检测的域名。",
]
```
* 出参
```
[
    	"ip":"ip",
    	"url":"ur;",
    	"ip_location":"ip映射地理位置",
    	"status":"状态",
    	"msg":"信息",
    	"redirect_count":"3xx跳转次数",
    	"http_code":"http状态码",
    	"dns_resolve_time":"解析时间 ms",
    	"http_conn_time":"连接时间 ms",
    	"http_pre_trans":"准备传输时间 ms",
    	"http_start_trans":"开始传输时间 ms",
    	"http_total_time":"总时间 ms",
    	"http_size_download":"下载大小 kb",
    	"http_speed_download":"下载速度 mb/s",
]
```

# 时间监控应用
## 背景说明
>监控服务器时间并实时上报消息到rabbit；消费端获取消息并进行时间判断是否正常。
## 启动说明
* 参照src/main/config.example.yml文件配置src/main/config.yml文件；注意interval_time参数
* 打包：go build src/main/index_monitor_time.go
```
//打包成64位linux可执行文件
env GOOS=linux GOARCH=amd64 go build src/main/index_monitor_time.go
//打包成64位windows可执行文件，并静默执行
env GOOS=windows GOARCH=amd64 go build src/main/index_monitor_time.go -ldflags "-H windowsgui"
```
## 执行
```
// linux
chmod 777 ***
nohup ./*** &
// windows
使用管理员权限运行.exe文件
```

# 域名监控应用
## 背景说明
>监控redis中配置的大量域名，并实时将异常域名数据上报消息到rabbit；消费端获取消息并进行报警服务。
## 启动
* 参照src/main/config.example.yml文件配置src/main/config.yml文件；注意interval_domain参数
* 打包：go build src/main/index_monitor_domain.go
```
//打包成64位linux可执行文件
env GOOS=linux GOARCH=amd64 go build src/main/index_monitor_domain.go
//打包成64位windows可执行文件，并静默执行
env GOOS=windows GOARCH=amd64 go build src/main/index_monitor_domain.go -ldflags "-H windowsgui"
```
## 执行
```
// linux
chmod 777 ***
nohup ./*** &
// windows
使用管理员权限运行.exe文件
```