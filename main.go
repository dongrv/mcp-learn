package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"
	"log"
	"mcp-learn/llm"
	"mcp-learn/utils"
	"os"
)

// 先编译MCP Server: mcp-servers/writer/build.bat

func init() {
	utils.RegisterMCPServer("writer", "./bin/writer.exe")
}

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY") // 环境变量
	if len(apiKey) == 0 {
		log.Fatalln("DEEPSEEK_API_KEY variable not set.")
	}

	client := llm.DeepSeekClient(apiKey)

	tools := []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "writer",
				Description: openai.String("Write a novel according to given message."),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"content": map[string]interface{}{
							"type":        "string",
							"description": "The content of novel.",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The name of novel.",
						},
					},
					"required": []string{"content", "name"},
				},
			},
		},
	}

	var messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("你需要调用 writer 工具，当有人希望你写一篇小说时"),
		openai.UserMessage("直接帮我生成一篇短篇推理悬疑小说，自行设定名字和内容，字数在800字左右，文笔风格参考《心理罪》，时代背景放在清末，内容探讨人性善恶。"),
	}

	log.Println("--- 开始请求LLM ---")
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:     "deepseek-chat",
		Messages:  messages,
		Tools:     tools,
		MaxTokens: openai.Int(5000),
	})
	if err != nil {
		log.Fatalln(err)
	}
	if len(resp.Choices) == 0 {
		log.Fatalln("没有返回可选择的聊天补全信息")
	}
	msg := resp.Choices[0].Message

	if msg.ToolCalls == nil {
		log.Fatalln("没有可调用的工具列表")
	}

	for _, tool := range msg.ToolCalls {
		log.Printf("调用工具的输入参数:%v\n", tool.Function.Arguments)
		functionName := tool.Function.Name
		var arguments map[string]interface{}
		if err := json.Unmarshal([]byte(tool.Function.Arguments), &arguments); err != nil {
			log.Fatalf("Failed to unmarshal function arguments: %v\n", err)
		}

		toolResult, err := utils.InvokeMCPTool(functionName, arguments)
		if err != nil {
			log.Printf("Tool call failed: %v\n", err)
			toolResult = fmt.Sprintf("Error executing tool: %v", err)
		}
		messages = append(messages, openai.SystemMessage(toolResult))
		log.Printf("--- Tool result: ---\n%s\n-------------------\n", toolResult)

		messages = append(messages, openai.ToolMessage(toolResult, tool.ID))
	}
	log.Println("--- Done ---")
}
