package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tuanta7/transcript/internal/core"
	"github.com/tuanta7/transcript/internal/ui"
	"github.com/tuanta7/transcript/pkg/audio"
	"github.com/tuanta7/transcript/pkg/transcriptor"
)

func main() {
	ctx := context.Background()

	mode := transcriptor.Mode(os.Getenv("TRANSCRIPTOR_MODE"))
	gc, err := transcriptor.NewClient(ctx, mode, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		fmt.Printf("Failed to create Gemini client: %v", err)
		os.Exit(1)
	}

	recorder := audio.NewRecorder()
	app := core.NewApplication(recorder, gc)

	model := ui.NewModel(app)
	_, err = tea.NewProgram(model).Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
