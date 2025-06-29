# MCP Demo Golang [第二课]

## 协议类型

- Tools
- Prompts
- Resources

## 传输方式

- stdio
- http
- sse

## Tools

每个`Tool`可以作为一个独立的`MCP Server`，但也可以在一个`MCP Server`中包含多个`Tool`，每个`Tool`对应一个`API`

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

## 参考

> https://modelcontextprotocol.io/quickstart/server
> https://modelcontextprotocol.io/quickstart/client

[上一课](https://github.com/guobinqiu/llm-chat)
[下一课](https://github.com/guobinqiu/mcp-host)

[Python版](https://github.com/guobinqiu/mcp-demo-python)
