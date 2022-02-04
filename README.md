# KnowlGraph

A publishing platform with user autonomy and shared community income where people can read and manage professional and in-depth knowledge on the graph.

### Build for linux

使用 ~~packr v2~~ `go:embed` 打包静态资源

```
set GOARCH=amd64
set GOOS=linux
go build --ldflags="-w -s" -o "knowlgraph.com"
upx -9 knowlgraph.com
echo "finished"
```

### Usage

```
chmod 775 knowlgraph.com
knowlgraph.com -config /your_path/config.json
```

### Test server response time

```shell
curl -w "\n   time_namelookup:  %{time_namelookup}\n      time_connect:  %{time_connect}\n   time_appconnect:  %{time_appconnect}\n     time_redirect:  %{time_redirect}\n  time_pretransfer:  %{time_pretransfer}\ntime_starttransfer:  %{time_starttransfer}\n                    ----------\n         time_total:  %{time_total}\n" http://localhost:20011/v1/words
```

### Deploy website

two modes:

1. Debug mode
2. Release mode

Debug 模式会打印具体的日志，包括数据库操作日志，且支持所有域的网络请求，前端 js 库是非压缩版的。
Release 模式则与之相反。
所以部署为 Release 模式时，请代理 Restful api 服务，否则存在端口号不一致而出现 CORS 访问异常。
