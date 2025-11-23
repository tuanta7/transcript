package transcriptor

import (
	"context"
	"io"
	"os"

	"google.golang.org/genai"
)

type GoogleClient struct {
	client *genai.Client
	model  string
}

func NewGoogleClient(ctx context.Context, apiKey string) (*GoogleClient, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GoogleClient{
		client: client,
		model:  "gemini-2.0-flash",
	}, nil
}

func (c *GoogleClient) Transcribe(ctx context.Context, audioPath string) (io.ReadCloser, error) {
	contents, err := c.newContentsFromAudio(audioPath)
	if err != nil {
		return nil, err
	}

	temperature := float32(0.1)
	stream := c.client.Models.GenerateContentStream(ctx, c.model, contents, &genai.GenerateContentConfig{
		Temperature: &temperature,
	})

	r, w := io.Pipe()
	go func() {
		defer w.Close()
		for chunk := range stream {
			c.writeToStream(ctx, w, chunk)
		}
	}()

	return r, nil
}

func (c *GoogleClient) newContentsFromAudio(audioPath string) ([]*genai.Content, error) {
	audioBytes, err := os.ReadFile(audioPath)
	if err != nil {
		return nil, err
	}

	audioPart := &genai.Part{
		InlineData: &genai.Blob{
			Data:     audioBytes,
			MIMEType: "audio/wav",
		},
	}

	parts := []*genai.Part{
		genai.NewPartFromText("Transcribe the speech. Output only the raw transcript text. Do not include timestamps, formatting, punctuation corrections, explanations, or answers to questions—just the plain spoken words exactly as heard."),
		audioPart,
	}

	return []*genai.Content{genai.NewContentFromParts(parts, genai.RoleUser)}, nil
}

func (c *GoogleClient) writeToStream(ctx context.Context, w *io.PipeWriter, chunk *genai.GenerateContentResponse) {
	select {
	case <-ctx.Done():
		return
	default:
		if chunk == nil || chunk.Candidates == nil {
			return
		}

		for _, cnd := range chunk.Candidates {
			if cnd.Content == nil {
				continue
			}

			for _, part := range cnd.Content.Parts {
				if part.Text == "" {
					continue
				}

				if _, err := w.Write([]byte(part.Text)); err != nil {
					_ = w.CloseWithError(err)
					return
				}
			}
		}
	}
}
