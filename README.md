# MCP Demo Golang

## tools 

协议类型

- stdio 标准输入输出
- http HTTP 请求/响应
- sse 服务端推流

### 运行 stdio

```
cd mcp-demo-golang
go build -o tools/stdio/ip_location_query/server/ip_location_query tools/stdio/ip_location_query/server/main.go
go run tools/stdio/ip_location_query/client/main.go tools/stdio/ip_location_query/server/ip_location_query
```

### 运行 http

```
cd mcp-demo-golang
go run tools/http/ip_location_query/server/main.go
go run tools/http/ip_location_query/client/main.go http://localhost:8080/mcp (in another terminal)
```

### 运行 sse

```
cd mcp-demo-golang
go run tools/sse/ip_location_query/server/main.go
go run tools/sse/ip_location_query/client/main.go http://localhost:8081/sse (in another terminal)
```

## prompts

```
cd mcp-demo-golang
go build -o prompts/stdio/ip_location_query/server/ip_location_query prompts/stdio/ip_location_query/server/main.go
go run prompts/stdio/ip_location_query/client/main.go prompts/stdio/ip_location_query/server/ip_location_query
```

## resources

```
cd mcp-demo-golang
go build -o resources/stdio/readme/server/readme resources/stdio/readme/server/main.go
go run resources/stdio/readme/client/main.go resources/stdio/readme/server/readme
```
