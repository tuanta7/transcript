# Transcript

A real-time desktop audio transcription application with a clean TUI interface. Transcribe your microphone audio on the fly using Google's Gemini API.

![Demo](demo.gif)

> Darts scene in Ted Lasso (2x speed)

## Quick Start

### Prerequisites

```bash
# Install required dependencies
./install.sh
```

This will install:

- `pulseaudio-utils` - For audio capture
- `ffmpeg` - For audio processing

## Configuration

Environment variables

| Variable            | Description           | Values            |
| ------------------- | --------------------- | ----------------- |
| `TRANSCRIPTOR_MODE` | Transcription backend | `google`, `local` |
| `GEMINI_API_KEY`    | Google Gemini API key | Your API key      |

### Get a Gemini API Key

- Visit [Google AI Studio](https://aistudio.google.com/app/api-keys)
- Create a new API key
- Add it to your `.env` file

## TODO

- [ ] Add support for local Gemma 3n model
- [ ] Support for multiple audio sources
- [ ] Export to text/markdown format
- [ ] Real-time word highlighting
- [ ] Custom recording duration settings
