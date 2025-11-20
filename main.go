package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tuanta7/transcript/pkg/gemini"
	"github.com/tuanta7/transcript/pkg/ui"
)

func main() {
	ctx := context.Background()

	gc, err := gemini.NewClient(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		fmt.Printf("Failed to create Gemini client: %v", err)
		os.Exit(1)
	}

	model := ui.NewModel(gc)

	_, err = tea.NewProgram(model).Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
