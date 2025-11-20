package ui

import "time"

type StopRecordingMsg struct{}

type RecordingDoneMsg struct {
	filename string
}

type TranscriptionDoneMsg struct {
	timestamp time.Time
	text      string
	err       error
}
