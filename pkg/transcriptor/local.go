package transcriptor

import (
	"context"
	"io"
	"net/http"
)

type LocalClient struct {
	httpClient             *http.Client
	cloudflareClientID     string
	cloudflareClientSecret string
}

func NewLocalClient(clientID, clientSecret string, httpClient *http.Client) (*LocalClient, error) {
	return &LocalClient{
		httpClient:             httpClient,
		cloudflareClientID:     clientID,
		cloudflareClientSecret: clientSecret,
	}, nil
}

func (l *LocalClient) Transcribe(ctx context.Context, audioPath string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://ollama.jodspace.work/o/generate", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("CF-Access-Client-Id", l.cloudflareClientID)
	request.Header.Add("CF-Access-Client-Secret", l.cloudflareClientSecret)
	request.Header.Add("Content-Type", "audio/wav")

	response, err := l.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// TODO: handle response
	return nil, nil
}
