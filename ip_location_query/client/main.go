package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// 创建 SSE transport
	sseTransport, err := transport.NewSSE("http://localhost:8080/sse")
	if err != nil {
		log.Fatalf("Failed to create SSE transport: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 启动 transport
	if err := sseTransport.Start(ctx); err != nil {
		log.Fatalf("Failed to start SSE transport: %v", err)
	}

	// 创建 MCP 客户端
	c := client.NewClient(sseTransport)
	defer c.Close()

	// 初始化 MCP 客户端
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "ip-location-client",
		Version: "1.0.0",
	}

	initResult, err := c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Printf("Connected to server: %s %s\n\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 构造调用 tool 请求参数
	req := mcp.CallToolRequest{}
	req.Params.Name = "ip_location_query"
	req.Params.Arguments = map[string]any{
		"ip": "183.193.158.228",
	}

	// 调用工具
	resp, err := c.CallTool(ctx, req)
	if err != nil {
		log.Fatalf("CallTool 调用失败: %v", err)
	}

	// 输出结果
	fmt.Printf("Tool 返回结果: %+v\n", resp)
}
