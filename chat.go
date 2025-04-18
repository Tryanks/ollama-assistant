package main

import (
	"bufio"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/openai/openai-go"
)

func ChatCompletion(c *fiber.Ctx) error {
	var chat OllamaChat
	err := c.BodyParser(&chat)
	if err != nil {
		log.Error(err)
		return err
	}

	params := openai.ChatCompletionNewParams{
		Model:    chat.Model,
		Messages: chat.GetOpenaiMessages(),
	}

	if !chat.Stream {
		return nonStreamingChatCompletion(c, chat.Model, params)
	}

	return streamingChatCompletion(c, chat.Model, params)
}

func streamingChatCompletion(c *fiber.Ctx, model string, params openai.ChatCompletionNewParams) error {
	stream := provider.Chat.Completions.NewStreaming(c.Context(), params)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		enc := sonic.ConfigDefault.NewEncoder(w)
		response := NewOllamaStreamResponse(model)

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

func nonStreamingChatCompletion(c *fiber.Ctx, model string, params openai.ChatCompletionNewParams) error {
	resp, err := provider.Chat.Completions.New(c.Context(), params)
	if err != nil {
		log.Error("Failed to create completion:", err)
		return err
	}

	response := NewOllamaStreamResponse(model)
	return c.JSON(response.End(resp.Choices[0].Message.Content))
}
