package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/tuanta7/ekko/internal/core"
	"github.com/tuanta7/ekko/pkg/logger"
)

type screen int

const (
	screenMenu screen = iota
	screenRecording
)

type Model struct {
	screen          screen
	cursor          int
	menuOptions     []string
	chunkDuration   time.Duration
	errorMsg        string
	sessionStopping bool

	spinner           spinner.Model
	transcript        viewport.Model
	transcriptContent string
	chunkCount        int
	sessionStart      time.Time

	app    *core.Application
	stream <-chan string
	logger *logger.FileLogger
}

func NewModel(app *core.Application) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	vp := viewport.New(100, 10)
	vp.SetContent("")

	return &Model{
		screen:        screenMenu,
		menuOptions:   []string{"Start Session", "Chunk Duration", "Exit"},
		spinner:       sp,
		transcript:    vp,
		app:           app,
		chunkDuration: 10 * time.Second,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		m.screen = screenRecording
		m.transcriptContent = ""
		m.transcript.SetContent("")
		m.transcript.YOffset = 0
		m.errorMsg = ""
		m.chunkCount = 0
		m.sessionStart = time.Now()

		var err error
		m.stream, err = m.app.Start(m.chunkDuration)
		if err != nil {
			return m, tea.Quit
		}

		return m, tea.Batch(m.spinner.Tick, m.waitForTranscript())
	case 1:
		// chunk duration, no action on selection
		return m, nil
	case len(m.menuOptions) - 1:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m *Model) handleKeyEvent(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenMenu:
		switch msg.String() {
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
		case "left":
			if m.cursor == 1 {
				if m.chunkDuration > time.Second {
					m.chunkDuration -= time.Second
				}
			}
		case "right":
			if m.cursor == 1 {
				if m.chunkDuration < 60*time.Second {
					m.chunkDuration += time.Second
				}
			}
		case "enter":
			return m.handleMenuSelection()
		}
	case screenRecording:
		switch msg.String() {
		case "ctrl+c", "q":
			_, _ = m.app.Stop()
			return m, tea.Quit
		case "s", "S":
			if m.sessionStopping {
				return m, nil // ignore spam
			}
			m.sessionStopping = true
			return m, m.sessionEnd()
		default:
			var cmd tea.Cmd
			m.transcript, cmd = m.transcript.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch mt := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyEvent(mt)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(mt)
		return m, cmd
	case transcriptChunkMsg:
		m.chunkCount++
		m.transcriptContent += mt.Text + "\n"
		wrapped := wordwrap.String(m.transcriptContent, m.transcript.Width-3)
		m.transcript.SetContent(wrapped)
		m.transcript.GotoBottom()
		return m, m.waitForTranscript()
	case sessionEndMsg:
		m.screen = screenMenu
		m.sessionStopping = false // reset guard
		if mt.Error != nil {
			m.errorMsg = fmt.Sprintf("Error: %v", mt.Error)
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	var b strings.Builder

	switch m.screen {
	case screenMenu:
		b.WriteString(subtitleStyle.Render(" Real-time Audio Transcription"))
		b.WriteString("\n")

		// Menu items
		var menuItems strings.Builder
		for i, choice := range m.menuOptions {
			icon := " "
			label := choice

			switch choice {
			case "Start Session":
				icon = "▶"
			case "Chunk Duration":
				icon = "⏱"
				durVal := durationValueStyle.Render(fmt.Sprintf("%ds", int(m.chunkDuration.Seconds())))
				label = fmt.Sprintf("Chunk Duration: %s  ◀ ▶", durVal)
			case "Exit":
				icon = "✕"
			}

			if m.cursor == i {
				cursor := cursorStyle.Render("●")
				menuItems.WriteString(fmt.Sprintf(" %s %s %s\n", cursor, icon, selectedStyle.Render(label)))
			} else {
				menuItems.WriteString(fmt.Sprintf("   %s %s\n", icon, normalStyle.Render(label)))
			}
		}
		b.WriteString(menuBoxStyle.Render(menuItems.String()))

		if m.errorMsg != "" {
			b.WriteString("\n")
			b.WriteString(errorStyle.Render(" ⚠ " + m.errorMsg + " "))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		help := fmt.Sprintf("%s navigate  %s adjust  %s select  %s quit",
			helpKeyStyle.Render("↑↓"),
			helpKeyStyle.Render("←→"),
			helpKeyStyle.Render("enter"),
			helpKeyStyle.Render("q"))
		b.WriteString(helpStyle.Render(help))

	case screenRecording:
		elapsed := time.Since(m.sessionStart).Round(time.Second)
		recDot := recordingDotStyle.Render("●")
		status := fmt.Sprintf("%s %s  Recording  •  %s elapsed  •  %d chunks",
			m.spinner.View(),
			recDot,
			elapsed.String(),
			m.chunkCount)
		b.WriteString(statusStyle.Render(status))
		b.WriteString("\n")
		b.WriteString(transcriptBoxStyle.Render(transcriptTextStyle.Render(m.transcript.View())))
		b.WriteString("\n\n")

		// Help
		help := fmt.Sprintf("%s scroll  %s stop & save  %s quit",
			helpKeyStyle.Render("↑↓"),
			helpKeyStyle.Render("s"),
			helpKeyStyle.Render("q"))
		b.WriteString(helpStyle.Render(help))
	}

	return b.String()
}
