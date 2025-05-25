package main

import (
	"context"
	"log"
	"os"

	_ "embed"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed README.md
var content string

func main() {
	// 确保日志输出到标准错误，避免破坏标准输出的 MCP JSON 通信
	log.SetOutput(os.Stderr)

	s := server.NewMCPServer(
		"docs",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	resource := mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// 返回内容时，绝对不要输出非 JSON 内容
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     content,
			},
		}, nil
	})

	// 启动 MCP 服务器，使用标准输入输出通信
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v", err)
	}
}
