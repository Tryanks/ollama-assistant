package main

import (
	"github.com/bytedance/sonic"
	"github.com/openai/openai-go"
	"time"
)

type Tags struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

type OllamaChat struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream    bool                   `json:"stream"`
	KeepAlive string                 `json:"keep_alive"`
	Format    sonic.NoCopyRawMessage `json:"format"`
	Options   struct {
		NumCtx int64 `json:"num_ctx"`
	} `json:"options"`
}

func (c *OllamaChat) GetOpenaiMessages() []openai.ChatCompletionMessageParamUnion {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(c.Messages))
	for k, message := range c.Messages {
		switch message.Role {
		case "system":
			messages[k] = openai.SystemMessage(message.Content)
		case "assistant":
			messages[k] = openai.AssistantMessage(message.Content)
		default:
			messages[k] = openai.UserMessage(message.Content)
		}
	}
	return messages
}

type OllamaStreamResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	Done               bool          `json:"done"`
	TotalDuration      int64         `json:"total_duration,omitempty"`
	LoadDuration       int64         `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64         `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       int64         `json:"eval_duration,omitempty"`
}

func NewOllamaStreamResponse(model string) OllamaStreamResponse {
	return OllamaStreamResponse{
		Model: model,
		Message: OllamaMessage{
			Role:    "assistant",
			Content: "",
		},
	}
}

func (c *OllamaStreamResponse) Next(char string) OllamaStreamResponse {
	c.CreatedAt = time.Now().Format(time.RFC3339)
	c.Message.Content = char
	return *c
}

func (c *OllamaStreamResponse) End(char string) OllamaStreamResponse {
	c.Done = true
	return c.Next(char)
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
