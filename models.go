package main

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"slices"
	"strings"
	"sync"
	"time"
)

// Global variables to store the model list and providers
var (
	cachedModelList Tags
	modelListMutex  sync.RWMutex
	providers       []OpenAIProvider
	providersMutex  sync.RWMutex
)

type Tags struct {
	Models []Model `json:"models"`
}

// ScanModels scans the available models from all providers and updates the cached model list
func ScanModels() error {
	// Get the list of providers
	newProviders := GetOpenAIProviders()

	// Create a new model list
	newTags := Tags{
		Models: []Model{},
	}

	// For each provider, get the list of models
	for i := range newProviders {
		provider := &newProviders[i]

		// Create a temporary OpenAI client for this provider
		client := openai.NewClient(
			option.WithBaseURL(provider.BaseURL),
			option.WithAPIKey(provider.APIKey),
		)

		// Get the list of models
		pages, err := client.Models.List(context.Background())
		if err != nil {
			log.Error("Failed to scan models from provider", provider.BaseURL, ":", err)
			continue
		}

		// Clear the provider's model list
		provider.Models = []string{}

		// Add each model to the provider's model list and the global model list
		for _, model := range pages.Data {
			if ModelBlockFilter.BlockString(model.ID) {
				continue
			}

			// Add to provider's model list
			provider.Models = append(provider.Models, model.ID)

			// Add to global model list
			newTags.Models = append(newTags.Models, Model{
				Name:  model.ID,
				Model: model.ID,
			})
		}

		log.Info("Found", len(provider.Models), "models from provider", provider.BaseURL)
	}

	// Sort the model list
	slices.SortFunc(newTags.Models, func(a, b Model) int {
		return strings.Compare(a.Name, b.Name)
	})

	// Update the cached model list and providers
	modelListMutex.Lock()
	cachedModelList = newTags
	modelListMutex.Unlock()

	providersMutex.Lock()
	providers = newProviders
	providersMutex.Unlock()

	log.Info("Model list updated, found", len(newTags.Models), "models from", len(newProviders), "providers")
	return nil
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
