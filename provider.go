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
func NewOpenAIChatProvider(model string) *OpenAIChatProvider {
	return &OpenAIChatProvider{
		model: model,
		provider: openai.NewClient(
			option.WithBaseURL(getEnv("API_BASE_URL")),
			option.WithAPIKey(getEnv("API_KEY")),
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
	return NewOpenAIChatProvider(model_name)
}
