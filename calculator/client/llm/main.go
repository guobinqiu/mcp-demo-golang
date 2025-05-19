package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	apiKey     = "sk-1d1bcf421aba4d998ff9c06bd39574c2"
	apiBaseURL = "https://api.deepseek.com/v1"
	model      = "deepseek-chat"
)

type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // raw json string
}

type ChatRequest struct {
	Model      string        `json:"model"`
	Messages   []ChatMessage `json:"messages"`
	Tools      []ToolSpec    `json:"tools,omitempty"`
	ToolChoice string        `json:"tool_choice,omitempty"` // "auto"
}

type ToolSpec struct {
	Type     string             `json:"type"`
	Function ToolFunctionSchema `json:"function"`
}

type ToolFunctionSchema struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  mcp.ToolInputSchema `json:"parameters"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

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

	// 转换为 LLM 所需格式
	var toolSpecs []ToolSpec
	for _, tool := range toolsResp.Tools {
		toolSpecs = append(toolSpecs, ToolSpec{
			Type: "function",
			Function: ToolFunctionSchema{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	// 用户输入
	userQuery := "请帮我计算 7 加 6 是多少"

	// 第一次请求 LLM
	messages := []ChatMessage{
		{Role: "user", Content: userQuery},
	}

	chatReq := ChatRequest{
		Model:      model,
		Messages:   messages,
		Tools:      toolSpecs,
		ToolChoice: "auto",
	}
	chatResp, err := sendChatRequest(chatReq)
	if err != nil {
		log.Fatalf("Chat request failed: %v", err)
	}

	reply := chatResp.Choices[0].Message

	// 是否需要调用工具
	if len(reply.ToolCalls) > 0 {
		for _, call := range reply.ToolCalls {
			toolName := call.Function.Name
			var args map[string]interface{}
			err := json.Unmarshal([]byte(call.Function.Arguments), &args)
			if err != nil {
				log.Fatalf("Invalid tool arguments: %v", err)
			}

			req := mcp.CallToolRequest{
				Request: mcp.Request{Method: "tools/call"},
			}
			req.Params.Name = toolName
			req.Params.Arguments = args

			result, err := c.CallTool(ctx, req)
			if err != nil {
				log.Fatalf("Tool call failed: %v", err)
			}

			// 添加 tool 响应消息
			toolReply := ChatMessage{
				Role:    "tool",
				Content: getToolContent(result),
			}
			messages = append(messages, ChatMessage{
				Role:      "assistant",
				ToolCalls: reply.ToolCalls,
			})
			messages = append(messages, toolReply)

			// 再次请求 LLM 获取最终回答
			finalReq := ChatRequest{
				Model:    model,
				Messages: messages,
			}
			finalResp, err := sendChatRequest(finalReq)
			if err != nil {
				log.Fatalf("Final chat request failed: %v", err)
			}

			fmt.Println("\n💬 最终回复：")
			fmt.Println(finalResp.Choices[0].Message.Content)
		}
	} else {
		fmt.Println("\n💬 模型直接回复：")
		fmt.Println(reply.Content)
	}
}

func sendChatRequest(req ChatRequest) (*ChatResponse, error) {
	jsonData, _ := json.Marshal(req)

	reqHttp, err := http.NewRequest("POST", apiBaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Authorization", "Bearer "+apiKey)
	reqHttp.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(reqHttp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("LLM error: %s", string(respBody))
	}

	var chatResp ChatResponse
	err = json.Unmarshal(respBody, &chatResp)
	return &chatResp, err
}

func getToolContent(result *mcp.CallToolResult) string {
	for _, c := range result.Content {
		if text, ok := c.(mcp.TextContent); ok {
			return text.Text
		}
	}
	return "[No text content returned from tool]"
}
