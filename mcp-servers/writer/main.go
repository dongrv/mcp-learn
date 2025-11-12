package main

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"os"
)

type Input struct {
	Name    string `json:"name"`    // 标题
	Content string `json:"content"` // 内容
}

type Output struct {
	Length int `json:"length"`
}

func Write(ctx context.Context, request *mcp.CallToolRequest, input Input) (
	result *mcp.CallToolResult,
	output Output,
	err error,
) {
	var file *os.File
	file, err = os.OpenFile("./novel/"+input.Name+".txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	output.Length, err = file.WriteString(input.Content)
	return
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "writer",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{Description: "write a novel", Name: "writer"}, Write)
	log.Println("The writer mcp server is running...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalln(err)
	}
}
