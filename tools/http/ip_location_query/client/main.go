package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建 HTTP transport
	httpTransport, err := transport.NewStreamableHTTP("http://localhost:8080/mcp")
	if err != nil {
		log.Fatalf("Failed to create HTTP transport: %v", err)
	}

	// 创建客户端实例，连接 MCP 服务端
	c := client.NewClient(httpTransport)
	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer c.Close()

	// 初始化 MCP 客户端（发送初始化请求，建立连接）
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

	// 列出所有可用工具
	fmt.Println("Available tools:")
	toolsResp, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range toolsResp.Tools {
		fmt.Println("name:", tool.Name)
		fmt.Println("description:", tool.Description)
		fmt.Println("parameters:", tool.InputSchema)
	}

	// 调用工具
	req := mcp.CallToolRequest{}
	req.Params.Name = "ip_location_query"
	req.Params.Arguments = map[string]any{
		"ip": "183.193.158.228",
	}
	res, err := c.CallTool(ctx, req)
	if err != nil {
		log.Fatalf("工具调用失败: %v", err)
	}

	// 输出结果
	fmt.Println("工具输出结果: ")
	printToolResult(res)
}

func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}
