package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

func (a *Application) StartSession() (<-chan string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isRunning {
		return nil, errors.New("session already running")
	}

	a.reset()
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.isRunning = true

	stream := make(chan string, 100)
	_ = os.Mkdir(".tmp", 0755)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(stream) // Signal consumers when done
		defer a.recover("transcribe")

		err := a.transcribe(stream)
		if err != nil && !errors.Is(err, context.Canceled) {
			stream <- err.Error()
		}
	}()

	go func() {
		defer wg.Done()
		defer a.recover("record")

		err := a.record(10 * time.Second)
		if err != nil && !errors.Is(err, context.Canceled) {
			stream <- err.Error()
		}
	}()

	go func() {
		wg.Wait()
		a.mu.Lock()
		a.isRunning = false
		a.mu.Unlock()
	}()

	return stream, nil
}

func (a *Application) StopSession() (string, error) {
	a.mu.Lock()
	if !a.isRunning {
		a.mu.Unlock()
		return "", errors.New("no active session")
	}

	if a.cancel != nil {
		a.cancel()
	}
	a.mu.Unlock()

	// graceful shutdown
	time.Sleep(100 * time.Millisecond)

	filename, err := a.saveSession()
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (a *Application) saveSession() (string, error) {
	a.session.mu.Lock()
	a.session.EndTime = time.Now()
	timestamp := a.session.StartTime.Format("20060102-150405")
	filename := fmt.Sprintf("transcript-%s.json", timestamp)

	data, err := json.MarshalIndent(a.session, "", "\t")
	a.session.mu.Unlock()

	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %w", err)
	}

	if err = os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	_ = os.RemoveAll(".tmp")
	return filename, nil
}

func (a *Application) recover(method string) {
	if r := recover(); r != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[%s] panic recovered: %v\n", method, r)
	}
}
