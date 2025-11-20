package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tuanta7/transcript/pkg/audio"
	"github.com/tuanta7/transcript/pkg/gemini"
)

type screen int

const (
	screenMenu screen = iota
	screenRecording
)

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginBottom(1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	resultStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

type Transcription struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Error     string    `json:"error,omitempty"`
}

type SessionData struct {
	StartTime      time.Time       `json:"start_time"`
	EndTime        time.Time       `json:"end_time"`
	Transcriptions []Transcription `json:"transcriptions"`
}

type Model struct {
	screen      screen
	cursor      int
	menuOptions []string
	spinner     spinner.Model

	recorder *audio.Recorder
	gemini   *gemini.Client

	fileQueue    []string
	recordingNum int
	currentFile  string

	recording    bool
	transcribing bool

	session *SessionData
}

func NewModel(geminiClient *gemini.Client) *Model {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)

	return &Model{
		screen:      screenMenu,
		cursor:      0,
		menuOptions: []string{"Start Recording", "Exit"},
		recorder:    audio.NewRecorder(),
		gemini:      geminiClient,
		spinner:     s,
		session:     &SessionData{},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch mt := msg.(type) {
	case tea.KeyMsg:
		switch m.screen {
		case screenMenu:
			switch mt.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.menuOptions)-1 {
					m.cursor++
				}
			case "enter":
				return m.handleMenuSelection()
			}

		case screenRecording:
			switch mt.String() {
			case "s", "S":
				return m, func() tea.Msg {
					return StopRecordingMsg{}
				}
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(mt)
		return m, cmd

	case RecordingDoneMsg:
		m.recording = false
		m.transcribing = true
		m.currentFile = mt.filename
		return m, m.transcribeAudio(mt.filename)

	case TranscriptionDoneMsg:
		m.transcribing = false

		result := Transcription{
			Timestamp: mt.timestamp,
		}

		if mt.err != nil {
			result.Error = mt.err.Error()
		} else {
			result.Text = mt.text
		}

		m.session.Transcriptions = append(m.session.Transcriptions, result)

		// Clean up audio file
		if m.currentFile != "" {
			_ = os.Remove(m.currentFile)
		}

		// Start next recording
		m.recordingNum++
		m.recording = true
		m.currentFile = fmt.Sprintf("output_%d.wav", m.recordingNum)
		return m, recordAudio(10*time.Second, m.currentFile)

	case StopRecordingMsg:
		// Save session to file
		m.saveSession()
		m.screen = screenMenu
		m.session.Transcriptions = []Transcription{}
		m.recordingNum = 0
		return m, nil
	}

	return m, nil
}

func (m *Model) View() string {
	var s strings.Builder

	switch m.screen {
	case screenMenu:
		s.WriteString(titleStyle.Render("🎤 Audio TranscriptionDoneMsg Tool"))
		s.WriteString("\n\n")
		s.WriteString("Select an option:\n\n")

		for i, choice := range m.menuOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				choice = selectedStyle.Render(choice)
			} else {
				choice = normalStyle.Render(choice)
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("↑/k: up • ↓/j: down • enter: select • q: quit"))

	case screenRecording:
		s.WriteString(titleStyle.Render("🔴 Recording Session"))
		s.WriteString("\n\n")

		if m.recording {
			s.WriteString(fmt.Sprintf("%s Recording audio...\n", m.spinner.View()))
		} else if m.transcribing {
			s.WriteString(fmt.Sprintf("%s Transcribing audio...\n", m.spinner.View()))
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press 's' to stop and save session"))
		s.WriteString("\n\n")

		if len(m.session.Transcriptions) > 0 {
			s.WriteString(resultStyle.Render("=== Transcriptions ==="))
			s.WriteString("\n\n")

			for _, result := range m.session.Transcriptions {
				timeStr := result.Timestamp.Format("15:04:05")
				s.WriteString(fmt.Sprintf("[%s] ", timeStr))

				if result.Error != "" {
					s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", result.Error)))
				} else {
					s.WriteString(result.Text)
				}
				s.WriteString("\n\n")
			}
		}
	}

	return s.String()
}

func (m *Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	if m.cursor == len(m.menuOptions)-1 {
		return m, tea.Quit
	}

	// Start Recording option
	m.screen = screenRecording
	m.session.StartTime = time.Now()
	m.recording = true
	m.recordingNum = 0
	m.currentFile = fmt.Sprintf("output_%d.wav", m.recordingNum)

	return m, tea.Batch(
		m.spinner.Tick,
		recordAudio(20*time.Second, m.currentFile),
	)
}

func (m *Model) saveSession() {
	timestamp := m.session.StartTime.Format("20060102-150405")
	filename := fmt.Sprintf("transcript-%s.json", timestamp)

	data, err := json.MarshalIndent(m.session, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(filename, data, 0644)
}

func recordAudio(duration time.Duration, outputFile string) tea.Cmd {
	return func() tea.Msg {
		recorder := audio.NewRecorder()
		err := recorder.Record(duration, outputFile)
		if err != nil {
			return TranscriptionDoneMsg{
				timestamp: time.Now(),
				text:      "",
				err:       fmt.Errorf("recording failed: %w", err),
			}
		}
		return RecordingDoneMsg{filename: outputFile}
	}
}

func (m *Model) transcribeAudio(audioFile string) tea.Cmd {
	return func() tea.Msg {
		timestamp := time.Now()
		ctx := context.Background()

		contents, err := m.gemini.NewContentsFromAudio(audioFile)
		if err != nil {
			return TranscriptionDoneMsg{
				text:      "",
				err:       fmt.Errorf("failed to create contents from audio: %w", err),
				timestamp: timestamp,
			}
		}

		text, err := m.gemini.Transcribe(ctx, contents)
		if err != nil {
			return TranscriptionDoneMsg{
				text:      "",
				err:       fmt.Errorf("failed to transcribe audio: %w", err),
				timestamp: timestamp,
			}
		}

		return TranscriptionDoneMsg{
			text:      text,
			timestamp: timestamp,
		}
	}
}
