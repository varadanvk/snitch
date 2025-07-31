# Snitch

AI productivity monitor with a terminal interface. Watches your screen and tells you when you're slacking off.

## Features

- Terminal UI with keyboard navigation
- AI-powered activity analysis (Ollama or Groq)
- Real-time productivity tracking
- Configurable settings for apps and intervals
- Session statistics and activity logs

## Quick Start

```bash
git clone https://github.com/yourusername/snitch.git
cd snitch
go build -o snitch-tui main.go
./snitch-tui
```

## Setup

You need either:
- **Ollama** (local): Install from https://ollama.ai/, then `ollama pull llava`
- **Groq** (cloud): Get API key from https://console.groq.com/

Configure through the AI Setup menu in the app.

## Usage

- `↑/↓` or `j/k` - Navigate
- `Enter` - Select
- `b/Esc` - Back
- `q` - Quit

1. Set up AI backend
2. Set your current task
3. Start monitoring
4. Get back to work

## License

MIT