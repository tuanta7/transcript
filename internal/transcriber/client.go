package transcriber

import (
	"context"
	"errors"
	"io"
)

type Mode string

const (
	WhisperMode Mode = "whisper"
	GeminiMode  Mode = "gemini"

	InitialPrompts string = "Transcribe the speech. Output only the raw transcript text. Do not include timestamps, formatting, punctuation corrections, explanations, or answers to questions—just the plain spoken words exactly as heard."
)

type Client interface {
	Transcribe(ctx context.Context, audioPath string) (io.ReadCloser, error)
	ResetContext(ctx context.Context) error
	Close() error
}

func NewClient(ctx context.Context, mode Mode, apiKey ...string) (Client, error) {
	switch mode {
	case WhisperMode:
		return NewLocalClient()
	case GeminiMode:
		if len(apiKey) != 1 {
			return nil, errors.New("invalid credentials")
		}
		return NewGeminiClient(ctx, apiKey[0])
	default:
		return nil, errors.New("invalid client mode, must be one of: local, google")
	}
}
