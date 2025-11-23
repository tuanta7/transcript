package transcriptor

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type Mode string

const (
	LocalMode  Mode = "local"
	GoogleMode Mode = "google"
)

type Client interface {
	Transcribe(ctx context.Context, audioPath string) (io.ReadCloser, error)
}

func NewClient(ctx context.Context, mode Mode, credentials ...string) (Client, error) {
	switch mode {
	case LocalMode:
		if len(credentials) != 2 {
			return nil, errors.New("invalid credentials")
		}
		return NewLocalClient(credentials[0], credentials[1], &http.Client{})
	case GoogleMode:
		if len(credentials) != 1 {
			return nil, errors.New("invalid credentials")
		}
		return NewGoogleClient(ctx, credentials[0])
	default:
		return nil, errors.New("invalid client mode, must be one of: local, google")
	}
}
