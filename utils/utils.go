package utils

import (
	"context"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"os/exec"
	"strings"
)

var registry = map[string]string{} // 注册表：mcp server工具名和对应调用路径

func RegisterMCPServer(typ string, cmd string) {
	registry[typ] = cmd
}

// InvokeMCPTool 连接并调用目标服务器的工具
func InvokeMCPTool(mcpToolName string, mcpArguments map[string]interface{}) (string, error) {
	cmd, ok := registry[strings.TrimSpace(mcpToolName)]
	if !ok {
		return "", fmt.Errorf("unknown tool alias: %s", mcpToolName)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "go-agent",
		Version: "v1.0.0",
	}, nil)

	cmdParts := strings.Fields(cmd)
	transport := &mcp.CommandTransport{Command: exec.Command(cmdParts[0], cmdParts[1:]...)}

	client.AddRoots(&mcp.Root{URI: "file://./"})

	ctx := context.Background()
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var resultText string
	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: mcpToolName, Arguments: mcpArguments})
	if err != nil {
		return "", err
	}
	if res.IsError {
		return "", fmt.Errorf("tool execution failed: %s", res.Content[0].(*mcp.TextContent).Text)
	}
	resultText = res.Content[0].(*mcp.TextContent).Text
	return resultText, nil
}
