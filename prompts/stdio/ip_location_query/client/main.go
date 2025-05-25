package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建客户端实例，连接 MCP 服务端
	c, err := client.NewStdioMCPClient(
		"bin/ip-location-server",
		[]string{},
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
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
	resp, err := c.CallTool(ctx, req)
	if err != nil {
		log.Fatalf("工具调用失败: %v", err)
	}

	// 调用Prompt
	toolResult, err := json.Marshal(resp.Content)
	if err != nil {
		log.Fatalf("Failed to marshal tool result: %v", err)
	}
	promptReq := mcp.GetPromptRequest{}
	promptReq.Params.Name = "ip_location_prompt"
	promptReq.Params.Arguments = map[string]string{
		"tool_result": string(toolResult),
	}
	result, err := c.GetPrompt(ctx, promptReq)
	if err != nil {
		log.Fatalf("GetPrompt error: %v", err)
	}

	// 输出结果
	fmt.Println(result)
}
