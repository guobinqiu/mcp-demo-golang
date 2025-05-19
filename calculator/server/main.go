package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Calculator Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add a calculator tool
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	// Add the calculator handler
	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Using helper functions for type-safe argument access
		// op, err := request.RequireString("operation")
		// if err != nil {
		// 	return mcp.NewToolResultError(err.Error()), nil
		// }
		op, ok := request.Params.Arguments["operation"].(string)
		if !ok {
			return nil, errors.New("operation must be a string")
		}

		// x, err := request.RequireFloat("x")
		// if err != nil {
		// 	return mcp.NewToolResultError(err.Error()), nil
		// }
		x, ok := request.Params.Arguments["x"].(float64)
		if !ok {
			return nil, errors.New("x must be a float")
		}

		// y, err := request.RequireFloat("y")
		// if err != nil {
		// 	return mcp.NewToolResultError(err.Error()), nil
		// }
		y, ok := request.Params.Arguments["y"].(float64)
		if !ok {
			return nil, errors.New("y must be a float")
		}

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = x / y
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
