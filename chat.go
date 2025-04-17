package main

import (
	"bufio"
	"context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/openai/openai-go"
	"time"
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
		resp, err := provider.Chat.Completions.New(context.Background(), params)
		if err != nil {
			log.Error("Failed to create completion:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate response",
			})
		}

		response := OllamaStreamResponse{
			Model:     chat.Model,
			CreatedAt: time.Now().Format(time.RFC3339),
			Message: OllamaMessage{
				Role:    "assistant",
				Content: resp.Choices[0].Message.Content,
			},
			Done: true,
		}

		return c.JSON(response)
	}

	c.Set("Content-Type", "application/x-ndjson")

	stream := provider.Chat.Completions.NewStreaming(context.Background(), params)
	acc := openai.ChatCompletionAccumulator{}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		enc := sonic.ConfigDefault.NewEncoder(w)
		resp := OllamaStreamResponse{
			Model: chat.Model,
			Message: OllamaMessage{
				Role:    "assistant",
				Content: "",
			},
		}
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			char := ""
			if len(chunk.Choices) > 0 {
				char = chunk.Choices[0].Delta.Content
			}

			err = enc.Encode(resp.Next(char))
			if err != nil {
				log.Warn(err)
			}
			_ = w.Flush()
		}

		err = enc.Encode(resp.End(""))
		if err != nil {
			log.Warn(err)
		}
		_ = w.Flush()
	})

	if stream.Err() != nil {
		log.Error(stream.Err())
	}

	return nil
}
