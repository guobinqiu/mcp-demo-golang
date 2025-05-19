package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// 该版本为保留上下文版本

// 有上下文（保留历史消息）：
// 你：谁是爱因斯坦？
// 🤖：爱因斯坦是20世纪著名的物理学家...
// 你：他最著名的理论是什么？
// 🤖：他最著名的理论是相对论，特别是广义相对论和狭义相对论。

// 无上下文（每轮都单独提问）：
// 你：谁是爱因斯坦？
// 🤖：爱因斯坦是20世纪著名的物理学家...
// 你：他最著名的理论是什么？
// 🤖：请明确你说的“他”是谁？
func main() {
	config := openai.DefaultConfig("sk-1d1bcf421aba4d998ff9c06bd39574c2")
	config.BaseURL = "https://api.deepseek.com/v1"

	client := openai.NewClientWithConfig(config)

	// 用于存储历史消息，实现多轮对话
	var messages []openai.ChatCompletionMessage

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("欢迎使用 Chat 模式，输入内容与模型对话，输入 `exit` 退出。")
	for {
		fmt.Print("\n你：")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "exit" || userInput == "quit" {
			break
		}
		if userInput == "" {
			continue
		}

		// 添加问题到历史消息
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "user", // 用户问
			Content: userInput,
		})

		response, err := chat(client, messages)
		if err != nil {
			fmt.Printf("请求失败: %v\n", err)
			continue
		}

		fmt.Printf("🤖：%s\n", response)

		// 添加回答到历史消息
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "assistant", //助手答
			Content: response,
		})
	}
}

func chat(client *openai.Client, messages []openai.ChatCompletionMessage) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    "deepseek-chat",
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}
	return resp.Choices[0].Message.Content, nil
}
