package core

import (
	"fmt"
	"time"

	"github.com/tuanta7/transcript/pkg/queue"
)

// record continuously records audio and enqueues for transcription
func (a *Application) record(duration time.Duration) error {
	source, err := a.recorder.GetSource(a.ctx)
	if err != nil {
		return fmt.Errorf("failed to get audio source: %w", err)
	}

	defer a.queue.Close()

	for {
		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		default:
			currentCount := a.counter.Add(1)
			fileName := fmt.Sprintf(".tmp/audio-%d.wav", currentCount)

			if err := a.recorder.Record(a.ctx, duration, source, fileName); err != nil {
				if a.ctx.Err() != nil {
					return a.ctx.Err()
				}
				return fmt.Errorf("recording failed: %w", err)
			}

			if err := a.queue.Enqueue(a.ctx, &queue.Message{
				Timestamp: time.Now(),
				FileName:  fileName,
			}); err != nil {
				if a.ctx.Err() != nil {
					return a.ctx.Err()
				}
				return err
			}
		}
	}
}
