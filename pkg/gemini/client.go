package gemini

import (
	"context"
	"os"
	"strings"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model  string
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		model:  "gemini-2.0-flash-lite",
	}, nil
}

func (c *Client) NewContentsFromAudio(audioPath string, prev ...string) ([]*genai.Content, error) {
	audioBytes, err := os.ReadFile(audioPath)
	if err != nil {
		return nil, err
	}

	audioPart := &genai.Part{
		InlineData: &genai.Blob{
			MIMEType: "audio/wav",
			Data:     audioBytes,
		},
	}

	var parts []*genai.Part

	if len(prev) > 0 {
		ctxText := strings.Join(prev, "\n\n---\n\n")
		parts = append(parts, genai.NewPartFromText("Previous transcript context:\n\n"+ctxText))
	}

	parts = append(parts,
		genai.NewPartFromText("Generate a transcript of the speech."),
		audioPart,
	)

	return []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}, nil
}

func (c *Client) Transcribe(ctx context.Context, contents []*genai.Content) (string, error) {
	temperature := float32(0.2)
	content, err := c.client.Models.GenerateContent(ctx, c.model, contents, &genai.GenerateContentConfig{
		Temperature: &temperature,
	})
	if err != nil {
		return "", err
	}

	return content.Text(), nil
}
