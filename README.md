# MCP Demo Golang

## tools 

协议类型

- stdio 标准输入输出
- http HTTP 请求/响应
- sse 服务端推流

### 运行 stdio

```
go build -o tools/stdio/ip_location_query/server/ip_location_query tools/stdio/ip_location_query/server/main.go
go run tools/stdio/ip_location_query/client/main.go tools/stdio/ip_location_query/server/ip_location_query
```

### 运行 http

```
go run tools/http/ip_location_query/server/main.go
go run tools/http/ip_location_query/client/main.go http://localhost:8080/mcp (in another terminal)
```

### 运行 sse

```
go run tools/sse/ip_location_query/server/main.go
go run tools/sse/ip_location_query/client/main.go http://localhost:8081/sse (in another terminal)
```

## prompts

## resources
