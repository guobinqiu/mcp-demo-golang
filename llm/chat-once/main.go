package main

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	config := openai.DefaultConfig("sk-1d1bcf421aba4d998ff9c06bd39574c2")
	config.BaseURL = "https://api.deepseek.com/v1"

	client := openai.NewClientWithConfig(config)

	result, err := chat(client, "解释量子计算")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}

func chat(client *openai.Client, content string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: content},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}
	return resp.Choices[0].Message.Content, nil
}
