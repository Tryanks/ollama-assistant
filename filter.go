package main

import "regexp"

type FilterCollection []*regexp.Regexp

var ModelBlockFilter FilterCollection = []*regexp.Regexp{
	// OpenAI models
	regexp.MustCompile("^dall-e-.*?"),               // Dall-E models
	regexp.MustCompile("^text-embedding-.*?"),       // Text Embeddings models
	regexp.MustCompile("^tts-.*?"),                  // TTS models
	regexp.MustCompile("^whisper-.*?"),              // Whisper models
	regexp.MustCompile("-\\d{4}-\\d{2}-\\d{2}$"),    // Date suffix
	regexp.MustCompile("-\\d{4}-\\d{2}-\\d{2}-ca$"), // Date suffix for azure
	regexp.MustCompile("-\\d{4}$"),                  // Short Date suffix
}

func (c FilterCollection) Block(b []byte) bool {
	for _, filter := range c {
		if filter.Match(b) {
			return true
		}
	}
	return false
}

func (c FilterCollection) BlockString(s string) bool {
	return c.Block([]byte(s))
}
