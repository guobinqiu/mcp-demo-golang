package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <COMMAND> [ARGS...]")
		os.Exit(1)
	}

	// 打开MCP服务器可执行文件
	cmd := exec.Command(os.Args[1])

	// 创建管道以捕获输出和输入
	stdout, _ := cmd.StdoutPipe()
	stdin, _ := cmd.StdinPipe()

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

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
