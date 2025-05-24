package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sashabaranov/go-openai"
)

type ChatClient struct {
	mcpClient    *client.Client
	openaiClient *openai.Client
	model        string
	messages     []openai.ChatCompletionMessage // 用于存储历史消息，实现多轮对话
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <COMMAND> [ARGS...]")
		os.Exit(1)
	}

	// 启动 MCP 客户端
	mcpClient, err := client.NewStdioMCPClient(
		os.Args[1],
		[]string{},
		os.Args[2:]...,
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer mcpClient.Close()

	// 初始化 MCP 客户端
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "ip-location-client",
		Version: "1.0.0",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Printf("Connected to server: %s %s\n\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_API_BASE")
	model := os.Getenv("OPENAI_API_MODEL")
	if apiKey == "" || baseURL == "" || model == "" {
		fmt.Println("检查环境变量设置")
		return
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	openaiClient := openai.NewClientWithConfig(config)

	cc := &ChatClient{
		mcpClient:    mcpClient,
		openaiClient: openaiClient,
		model:        model,
		messages:     make([]openai.ChatCompletionMessage, 0),
	}
	cc.ChatLoop()
}

func (cc *ChatClient) ChatLoop() {
	fmt.Print("Type your queries or 'quit' to exit.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nUser: ")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if strings.ToLower(userInput) == "quit" {
			break
		}
		if userInput == "" {
			continue
		}

		response, err := cc.ProcessQuery(userInput)
		if err != nil {
			fmt.Printf("请求失败: %v\n", err)
			continue
		}

		fmt.Printf("Assistant: %s\n", response)
	}
}

func (cc *ChatClient) ProcessQuery(userInput string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 列出所有可用工具
	availableTools := []openai.Tool{}
	toolsResp, err := cc.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
	}
	for _, tool := range toolsResp.Tools {
		// fmt.Println("name:", tool.Name)
		// fmt.Println("description:", tool.Description)
		// fmt.Println("parameters:", tool.InputSchema)
		availableTools = append(availableTools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	finalText := []string{}

	// 首轮交互
	cc.messages = append(cc.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	resp, err := cc.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    cc.model,
		Messages: cc.messages,
		Tools:    availableTools,
	})
	if err != nil {
		return "", err
	}

	for _, choice := range resp.Choices {
		message := choice.Message
		if message.Content != "" { // 若直接生成文本
			finalText = append(finalText, message.Content)
		} else if len(message.ToolCalls) > 0 { // 若调用工具
			toolCallMessages := []openai.ChatCompletionMessage{}

			for _, toolCall := range message.ToolCalls {
				toolName := toolCall.Function.Name
				toolArgsRaw := toolCall.Function.Arguments
				// fmt.Println("=====toolCall.Function.Arguments:", toolArgsRaw)
				var toolArgs map[string]interface{}
				_ = json.Unmarshal([]byte(toolArgsRaw), &toolArgs)

				// 调用工具
				req := mcp.CallToolRequest{}
				req.Params.Name = toolName
				req.Params.Arguments = toolArgs
				resp, err := cc.mcpClient.CallTool(ctx, req)
				if err != nil {
					log.Printf("工具调用失败: %v", err)
					continue
				}

				// 构造 tool message
				toolCallMessages = append(toolCallMessages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					ToolCallID: toolCall.ID,
					Content:    fmt.Sprintf("%s", resp.Content),
				})
			}

			// 添加 assistant tool call 信息
			cc.messages = append(cc.messages, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				Content:   "",
				ToolCalls: message.ToolCalls,
			})

			// 添加 tool 响应
			cc.messages = append(cc.messages, toolCallMessages...)

			// debug
			// b, _ := json.MarshalIndent(cc.messages, "", "  ")
			// fmt.Println("Sending messages to OpenAI:\n", string(b))

			// 再次发送给模型
			nextResponse, err := cc.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
				Model:    cc.model,
				Messages: cc.messages,
			})
			if err != nil {
				return "", err
			}

			for _, nextChoice := range nextResponse.Choices {
				if nextChoice.Message.Content != "" {
					finalText = append(finalText, nextChoice.Message.Content)
				}
			}
		}
	}
	return strings.Join(finalText, "\n"), nil
}
