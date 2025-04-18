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
	Tools     []Tool                 `json:"tools"`
	Options   struct {
		NumCtx int64 `json:"num_ctx"`
	} `json:"options"`
}

func BindOllamaChat(binder func(any) error) (chat OllamaChat, err error) {
	err = binder(&chat)
	return
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

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string                            `json:"type"`
			Properties map[string]sonic.NoCopyRawMessage `json:"properties"`
			Required   []string                          `json:"required"`
		} `json:"parameters"`
	} `json:"function"`
}

type OllamaStreamResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	DoneReason         string        `json:"done_reason,omitempty"`
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
	if c.DoneReason != "stop" {
		c.Message.ToolCalls = nil
	}
	return *c
}

func (c *OllamaStreamResponse) End(char string) OllamaStreamResponse {
	c.Done = true
	return c.Next(char)
}

func (c *OllamaStreamResponse) Call() OllamaStreamResponse {
	c.DoneReason = "stop"
	return c.End("")
}

type OllamaMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	ToolCalls []struct {
		Function ToolCall `json:"function"`
	} `json:"tool_calls,omitempty"`
}

func (m *OllamaMessage) AddToolCall(call ToolCall) {
	m.ToolCalls = append(m.ToolCalls, struct {
		Function ToolCall `json:"function"`
	}{Function: call})
}

type ToolCall struct {
	Name      string                            `json:"name"`
	Arguments map[string]sonic.NoCopyRawMessage `json:"arguments"`
}
