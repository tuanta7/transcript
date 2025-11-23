package core

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tuanta7/transcript/pkg/audio"
	"github.com/tuanta7/transcript/pkg/queue"
	"github.com/tuanta7/transcript/pkg/transcriptor"
)

type Transcription struct {
	Timestamp time.Time `json:"timestamp"`
	FileName  string    `json:"filename"`
	Text      string    `json:"text"`
	Error     error     `json:"error,omitempty"`
}

type Session struct {
	mu             sync.Mutex
	StartTime      time.Time       `json:"start_time"`
	EndTime        time.Time       `json:"end_time"`
	Transcriptions []Transcription `json:"transcriptions"`
}

type Application struct {
	mu        sync.Mutex
	isRunning bool
	ctx       context.Context
	cancel    context.CancelFunc

	session  *Session
	queue    *queue.RecordQueue
	counter  atomic.Uint32
	recorder *audio.Recorder
	gemini   transcriptor.Client
}

func NewApplication(recorder *audio.Recorder, gemini transcriptor.Client) *Application {
	return &Application{
		recorder: recorder,
		gemini:   gemini,
		session:  &Session{},
	}
}

func (a *Application) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isRunning
}

func (a *Application) reset() {
	a.session = &Session{
		StartTime:      time.Now(),
		Transcriptions: make([]Transcription, 0),
	}
	a.queue = queue.NewRecordQueue()
	a.counter.Store(0)
}
