package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"ip-location-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Add an ip tool
	ipTool := mcp.NewTool("ip_location_query",
		mcp.WithDescription("查询IP地址的地理位置"),
		mcp.WithString("ip",
			mcp.Required(),
			mcp.Description("要查询的IP地址"),
		),
	)

	// Add the tool handler
	s.AddTool(ipTool, ipQueryHandler)

	// Add an ip prompt
	prompt := mcp.NewPrompt("ip_location_prompt",
		mcp.WithPromptDescription("格式化IP地址查询结果"),
		mcp.WithArgument("tool_result",
			mcp.ArgumentDescription("IP地址查询结果"),
			mcp.RequiredArgument(),
		),
	)

	// Add the prompt handler
	s.AddPrompt(prompt, promptHandler)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func ipQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ip, ok := request.GetArguments()["ip"].(string)
	if !ok {
		return nil, errors.New("ip must be a string")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("无效的IP地址")
	}

	// 调用外部 IP 地理位置服务
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体错误: %v", err)
	}

	// return &mcp.CallToolResult{
	// 	Content: []mcp.Content{
	// 		mcp.TextContent{
	// 			Type: "text",
	// 			Text: string(data),
	// 		},
	// 	},
	// }, nil
	return mcp.NewToolResultText(string(data)), nil
}

func promptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "格式化后的IP地址查询结果",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(`你将接收到来自工具的响应JSON，它包含country、regionName、isp等字段。
						请将这个JSON转换为简洁的中文自然语言句子。
						输出格式：IP所在地为<country>-<regionName>，网络服务商为:<isp>。
						JSON:%s`, request.Params.Arguments["tool_result"]),
				},
			},
		},
	}, nil
}
