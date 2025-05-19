package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"Demo",
		"1.0.0",
		server.WithToolCapabilities(false),
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
		content, err := os.ReadFile("README.md")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})

	// 注册一个 Prompt：构造 SQL 查询
	s.AddPrompt(mcp.NewPrompt("query_builder",
		mcp.WithPromptDescription("SQL query builder assistance"),
		mcp.WithArgument("table",
			mcp.ArgumentDescription("Table to generate SQL for"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		table := req.Params.Arguments["table"]
		return mcp.NewGetPromptResult(
			"SQL query builder",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf("Generate a SQL query for table: %s", table)),
				),
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent("Based on the following schema and guidance:")),
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewEmbeddedResource(mcp.BlobResourceContents{
						URI:      "docs://readme",
						MIMEType: "text/markdown",
					}),
				),
			},
		), nil
	})

	// 启动服务器
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
