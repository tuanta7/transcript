package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tuanta7/transcript/internal/core"
)

type screen int

const (
	screenMenu screen = iota
	screenRecording
)

type Model struct {
	screen      screen
	cursor      int
	menuOptions []string

	spinner    spinner.Model
	transcript textarea.Model

	app    *core.Application
	stream <-chan string
}

func NewModel(app *core.Application) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	ta := textarea.New()
	ta.Placeholder = "Transcription will appear here..."
	ta.ShowLineNumbers = true
	ta.SetWidth(90)
	ta.SetHeight(10)
	ta.KeyMap.DeleteCharacterBackward.SetEnabled(false)
	ta.KeyMap.DeleteCharacterForward.SetEnabled(false)
	ta.KeyMap.DeleteWordBackward.SetEnabled(false)
	ta.KeyMap.DeleteWordForward.SetEnabled(false)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &Model{
		screen:      screenMenu,
		menuOptions: []string{"Start Session", "Exit"},
		spinner:     sp,
		transcript:  ta,
		app:         app,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		m.screen = screenRecording
		m.transcript.Reset()

		var err error
		m.stream, err = m.app.StartSession()
		if err != nil {
			return m, tea.Quit
		}

		return m, tea.Batch(m.spinner.Tick, m.waitForTranscript())
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
		case "enter":
			return m.handleMenuSelection()
		}
	case screenRecording:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.app.IsRunning() {
				_, _ = m.app.StopSession()
			}
			return m, tea.Quit
		case "s", "S":
			filename, err := m.app.StopSession()
			return m, m.sessionEnd(filename, err)
		default:
			// Pass other keys to textarea for scrolling
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
		m.transcript.SetValue(m.transcript.Value() + mt.Text)
		m.transcript.CursorEnd()
		return m, m.waitForTranscript()
	case sessionEndMsg:
		m.screen = screenMenu
		return m, nil
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	var b strings.Builder

	switch m.screen {
	case screenMenu:
		b.WriteString(titleStyle.Render("Audio Transcription"))
		b.WriteString("\n")

		for i, choice := range m.menuOptions {
			cursor := " "
			label := choice
			if m.cursor == i {
				cursor = ">"
				label = selectedStyle.Render(choice)
			} else {
				label = normalStyle.Render(choice)
			}
			b.WriteString(fmt.Sprintf("%s %s\n", cursor, label))
		}

		b.WriteString(helpStyle.Render("up/down: navigate • enter: select • q: quit"))
	case screenRecording:
		b.WriteString(titleStyle.Render("Recording Session"))
		b.WriteString("\n")
		b.WriteString(statusStyle.Render(m.spinner.View() + " Recording..."))
		b.WriteString("\n")
		b.WriteString(transcriptBoxStyle.Render(m.transcript.View()))
		b.WriteString(helpStyle.Render("↑/↓: scroll • s: stop and save • q: quit"))
	}

	return b.String()
}
