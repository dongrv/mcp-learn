package llm

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func DeepSeekClient(apiKey string) *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://api.deepseek.com/v1"),
	)
	return &client
}
