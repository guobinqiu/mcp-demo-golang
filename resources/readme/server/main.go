package main

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed README.md
var content string

func main() {
	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	// 创建资源
	resource := mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)

	// 添加资源处理器
	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		dir, _ := os.Getwd()
		fmt.Println("Current working dir:", dir)

		// content, err := os.ReadFile("/Users/guobin/workspace/mcp-demo-golang/readme/server/README.md")
		// if err != nil {
		// 	return nil, err
		// }

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				//Text:     string(content),
				Text: content,
			},
		}, nil
	})

	// 启动服务器
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
