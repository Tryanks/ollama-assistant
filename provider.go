package main

import (
	"bufio"
	"bytes"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatProvider interface {
	Model() string
	ChatCompletion(*fiber.Ctx, OllamaChat) error
}

// OpenAIChatProvider implements the ChatProvider interface for OpenAI
type OpenAIChatProvider struct {
	model    string
	provider openai.Client
}

// NewOpenAIChatProvider creates a new OpenAI chat provider
func NewOpenAIChatProvider(model string, baseURL string, apiKey string) *OpenAIChatProvider {
	return &OpenAIChatProvider{
		model: model,
		provider: openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(apiKey),
		),
	}
}

// Model returns the model name
func (p *OpenAIChatProvider) Model() string {
	return p.model
}

// ChatCompletion handles chat completion requests
func (p *OpenAIChatProvider) ChatCompletion(c *fiber.Ctx, chat OllamaChat) error {
	// Create OpenAI request
	params := openai.ChatCompletionNewParams{
		Model:    p.model,
		Messages: chat.GetOpenaiMessages(),
	}

	// Tool calling
	if chat.Tools != nil && !chat.Stream {
		// TODO:            ^ I don't know how to compatible this with Streaming
		params.Tools = make([]openai.ChatCompletionToolParam, len(chat.Tools))
		for k, tool := range chat.Tools {
			params.Tools[k] = openai.ChatCompletionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        tool.Function.Name,
					Description: openai.String(tool.Function.Description),
					Parameters: openai.FunctionParameters{
						"type":       tool.Function.Parameters.Type,
						"properties": tool.Function.Parameters.Properties,
						"required":   tool.Function.Parameters.Required,
					},
				},
			}
		}
	}

	// Structured outputs
	if chat.Format != nil {
		if bytes.Equal(chat.Format, []byte{'j', 's', 'o', 'n'}) {
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfText: &openai.ResponseFormatTextParam{
					Type: "json_object",
				},
			}
		} else {
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:   "ollama_chat",
						Schema: chat.Format,
					},
				},
			}
		}
	}

	// Streaming responses
	if !chat.Stream {
		return p.nonStreamingChatCompletion(c, params)
	}

	return p.streamingChatCompletion(c, params)
}

func (p *OpenAIChatProvider) streamingChatCompletion(c *fiber.Ctx, params openai.ChatCompletionNewParams) error {
	stream := p.provider.Chat.Completions.NewStreaming(c.Context(), params)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		enc := sonic.ConfigDefault.NewEncoder(w)
		response := NewOllamaStreamResponse(p.model)

		for stream.Next() {
			chunk := stream.Current()

			char := ""
			if len(chunk.Choices) > 0 {
				char = chunk.Choices[0].Delta.Content
			}

			err := enc.Encode(response.Next(char))
			if err != nil {
				log.Warn(err)
			}
			_ = w.Flush()
		}

		err := enc.Encode(response.End(""))
		if err != nil {
			log.Warn(err)
		}
		_ = w.Flush()
	})

	err := stream.Err()
	if err != nil {
		log.Error(err)
	}

	return err
}

func (p *OpenAIChatProvider) nonStreamingChatCompletion(c *fiber.Ctx, params openai.ChatCompletionNewParams) error {
	resp, err := p.provider.Chat.Completions.New(c.Context(), params)
	if err != nil {
		log.Error("Failed to create completion:", err)
		return err
	}

	response := NewOllamaStreamResponse(p.model)
	toolCalls := resp.Choices[0].Message.ToolCalls
	if toolCalls != nil {
		for _, toolCall := range toolCalls {
			tool := ToolCall{}
			tool.Name = toolCall.Function.Name
			err := sonic.UnmarshalString(toolCall.Function.Arguments, tool.Arguments)
			if err != nil {
				log.Error(err)
				continue
			}
			response.Message.AddToolCall(tool)
		}
		if response.Message.ToolCalls != nil {
			return c.JSON(response.Call())
		}
	}
	return c.JSON(response.End(resp.Choices[0].Message.Content))
}

// GetProvider returns a ChatProvider for the given model name
func GetProvider(model_name string) ChatProvider {
	// Lock the providers mutex to ensure thread safety
	providersMutex.RLock()
	defer providersMutex.RUnlock()

	// If no providers are available, initialize them
	if len(providers) == 0 {
		// Release the read lock before calling ScanModels
		providersMutex.RUnlock()
		ScanModels()
		// Re-acquire the read lock
		providersMutex.RLock()

		// If still no providers, use default
		if len(providers) == 0 {
			return NewOpenAIChatProvider(model_name, getEnv("API_BASE_URL"), getEnv("API_KEY"))
		}
	}

	// Find the provider that supports this model
	for _, provider := range providers {
		for _, model := range provider.Models {
			if model == model_name {
				return NewOpenAIChatProvider(model_name, provider.BaseURL, provider.APIKey)
			}
		}
	}

	// If no provider supports this model, use the first provider
	return NewOpenAIChatProvider(model_name, providers[0].BaseURL, providers[0].APIKey)
}
