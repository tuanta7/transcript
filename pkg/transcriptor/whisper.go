package transcriptor

import (
	"context"
	"io"
)

type WhisperClient struct {
}

func NewLocalClient() *WhisperClient {
	return &WhisperClient{}
}

func (l *WhisperClient) Transcribe(ctx context.Context, audioPath string) (io.ReadCloser, error) {
	return nil, nil
}
