package main

import (
	"context"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		client: client,
		model:  "gemini-2.0-flash-lite",
	}, nil
}

func (c *GeminiClient) Transcribe(ctx context.Context) (string, error) {
	temperature := float32(0.2)
	content, err := c.client.Models.GenerateContent(ctx, c.model, nil, &genai.GenerateContentConfig{
		Temperature: &temperature,
	})
	if err != nil {
		return "", err
	}

	return content.Text(), nil
}
