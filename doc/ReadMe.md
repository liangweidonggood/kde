# 架构
go1.16
go mod init kde

# 数据

# 启动
```
go run main.go -c config.yaml
```
# 设备指令
```
http://localhost:8089/api/v1/deviceCmd
Post
{
    "sn":"9031",
    "param":"EbyteNETAT+SOCK"
}
AT*GPSCFG=0,5#
AT*REGPKG=0#
查询第一路：EbyteNETAT+SOCK
查询第二路：EbyteNETAT+SOCK1
查询第三路：EbyteNETAT+SOCK2
重启：EbyteNETAT+REBT
配置第三路：EbyteNETAT+SOCK2=1,TCPC,zl.glodon.com,2057

```
