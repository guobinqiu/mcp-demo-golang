package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	// 打开MCP服务器可执行文件
	cmd := exec.Command("../server/readme") // 替换为你的MCP服务器可执行文件路径

	// 创建管道以捕获输出和输入
	stdout, _ := cmd.StdoutPipe()
	stdin, _ := cmd.StdinPipe()

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	/*
		MCP 协议中所有标准方法（methods）的常量定义

		const (
			// MethodInitialize initiates connection and negotiates protocol capabilities.
			// https://modelcontextprotocol.io/specification/2024-11-05/basic/lifecycle/#initialization
			MethodInitialize MCPMethod = "initialize"

			// MethodPing verifies connection liveness between client and server.
			// https://modelcontextprotocol.io/specification/2024-11-05/basic/utilities/ping/
			MethodPing MCPMethod = "ping"

			// MethodResourcesList lists all available server resources.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/resources/
			MethodResourcesList MCPMethod = "resources/list"

			// MethodResourcesTemplatesList provides URI templates for constructing resource URIs.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/resources/
			MethodResourcesTemplatesList MCPMethod = "resources/templates/list"

			// MethodResourcesRead retrieves content of a specific resource by URI.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/resources/
			MethodResourcesRead MCPMethod = "resources/read"

			// MethodPromptsList lists all available prompt templates.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/prompts/
			MethodPromptsList MCPMethod = "prompts/list"

			// MethodPromptsGet retrieves a specific prompt template with filled parameters.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/prompts/
			MethodPromptsGet MCPMethod = "prompts/get"

			// MethodToolsList lists all available executable tools.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/tools/
			MethodToolsList MCPMethod = "tools/list"

			// MethodToolsCall invokes a specific tool with provided parameters.
			// https://modelcontextprotocol.io/specification/2024-11-05/server/tools/
			MethodToolsCall MCPMethod = "tools/call"

			// MethodSetLogLevel configures the minimum log level for client
			// https://modelcontextprotocol.io/specification/2025-03-26/server/utilities/logging
			MethodSetLogLevel MCPMethod = "logging/setLevel"

			// MethodNotificationResourcesListChanged notifies when the list of available resources changes.
			// https://modelcontextprotocol.io/specification/2025-03-26/server/resources#list-changed-notification
			MethodNotificationResourcesListChanged = "notifications/resources/list_changed"

			MethodNotificationResourceUpdated = "notifications/resources/updated"

			// MethodNotificationPromptsListChanged notifies when the list of available prompt templates changes.
			// https://modelcontextprotocol.io/specification/2025-03-26/server/prompts#list-changed-notification
			MethodNotificationPromptsListChanged = "notifications/prompts/list_changed"

			// MethodNotificationToolsListChanged notifies when the list of available tools changes.
			// https://spec.modelcontextprotocol.io/specification/2024-11-05/server/tools/list_changed/
			MethodNotificationToolsListChanged = "notifications/tools/list_changed"
		)
	*/
	// 写入请求到stdin
	_, err := stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"docs://readme"}}` + "\n"))
	if err != nil {
		fmt.Println("Error writing to stdin:", err)
		return
	}
	stdin.Close()

	// 读取stdout中的响应
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println("Response from server:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}

	// 等待命令完成
	cmd.Wait()
}
