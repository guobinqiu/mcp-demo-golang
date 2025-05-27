# MCP Demo Golang [第二课]

[第一课](https://github.com/guobinqiu/llm-chat)

[第三课](https://github.com/guobinqiu/mcp-host)

## 协议类型

- Tools
- Prompts
- Resources

## 传输方式

- stdio
- http
- sse

## Tools

### stdio

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

### http

```
cd mcp-demo-golang
go run tools/http/ip_location_query/server/main.go
go run tools/http/ip_location_query/client/main.go (in another terminal)
```

### sse

```
cd mcp-demo-golang
go run tools/sse/ip_location_query/server/main.go
go run tools/sse/ip_location_query/client/main.go (in another terminal)
```

## Prompts

```
cd mcp-demo-golang
go build -o bin/ip-location-server prompts/stdio/ip_location_query/server/main.go
go run prompts/stdio/ip_location_query/client/main.go
```

## Resources

```
cd mcp-demo-golang
go build -o bin/docs-server resources/stdio/docs/server/main.go
go run resources/stdio/docs/client/main.go
```

## Python版

> https://github.com/guobinqiu/mcp-demo-python
