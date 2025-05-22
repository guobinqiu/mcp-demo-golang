package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <COMMAND> [ARGS...]")
		os.Exit(1)
	}

	// 启动 MCP 客户端
	c, err := client.NewStdioMCPClient(
		os.Args[1],
		[]string{},
		os.Args[2:]...,
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// 初始化 MCP 客户端
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "calculator-client",
		Version: "1.0.0",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	// 输出结果
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
