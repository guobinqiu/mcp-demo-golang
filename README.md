# MCP Demo Golang

## tools 

协议类型

- stdio 标准输入输出
- http HTTP 请求/响应
- sse 服务端推流

### 运行 stdio

ip_location_query

```
cd mcp-demo-golang
go build -o bin/ip-location-server tools/stdio/ip_location_query/server/main.go
go run tools/stdio/ip_location_query/client/main.go
```

ip_location_query with LLM

```
cd mcp-demo-golang
go build -o bin/ip-location-server tools/stdio/ip_location_query/server/main.go
go run tools/stdio/ip_location_query/llm-client/main.go
```

calculator

```
cd mcp-demo-golang
go build -o bin/calculator-server tools/stdio/calculator/server/main.go
go run tools/stdio/calculator/client/main.go
```

calculator with LLM

```
cd mcp-demo-golang
go build -o bin/calculator-server tools/stdio/calculator/server/main.go
go run tools/stdio/calculator/llm-client/main.go
```

### 运行 http

```
cd mcp-demo-golang
go run tools/http/ip_location_query/server/main.go
go run tools/http/ip_location_query/client/main.go (in another terminal)
```

### 运行 sse

```
cd mcp-demo-golang
go run tools/sse/ip_location_query/server/main.go
go run tools/sse/ip_location_query/client/main.go (in another terminal)
```

## prompts

```
cd mcp-demo-golang
go build -o bin/ip-location-server prompts/stdio/ip_location_query/server/main.go
go run prompts/stdio/ip_location_query/client/main.go
```

## resources

```
cd mcp-demo-golang
go build -o bin/docs-server resources/stdio/docs/server/main.go
go run resources/stdio/docs/client/main.go
```
