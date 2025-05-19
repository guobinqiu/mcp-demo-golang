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
	c, err := client.NewStdioMCPClient(
		"/Users/guobin/workspace/mcp-demo-golang/calculator/server/calculator",
		[]string{},
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 初始化 MCP 客户端
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Calculator Client",
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
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}
	fmt.Println()

	// 调用 calculate 工具，计算 5 + 4
	req := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	req.Params.Name = "calculate"
	req.Params.Arguments = map[string]any{
		"operation": "add",
		"x":         5,
		"y":         4,
	}

	res, err := c.CallTool(ctx, req)
	if err != nil {
		log.Fatalf("Failed to call calculate tool: %v", err)
	}

	fmt.Println("Calculate tool result:")
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
