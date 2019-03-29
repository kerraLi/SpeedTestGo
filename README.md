# WEB应用
## 配置管理

* 默认端口
>8888
* 作用
>对外提供配置文件上传至服务器接口，统一服务器配置管理、操作；统一管理繁杂的服务器配置。

### API：配置上传

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


## 单域名速度检测
* 默认端口
>2019
* 作用
>服务器作为测速节点，对外提供域名速度检测接口，快速检测域名再当前服务器区域解析结果、速度等。

### API：速度检测
