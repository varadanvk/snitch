# Snitch

AI productivity monitor with a terminal interface. Watches your screen and tells you when you're slacking off. I like to call it a "lock-in" manager

## Features

- Terminal UI with keyboard navigation
- AI-powered activity analysis (Ollama or Groq)
- Real-time productivity tracking
- Configurable settings for apps and intervals
- Session statistics and activity logs

## Installation

### Option 1: Install via Go (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/varadanvk/snitch/main/install.sh | bash
```

### Option 2: Manual Install

```bash
go install github.com/varadanvk/snitch@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/varadanvk/snitch.git
cd snitch
make build
./snitch
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

1. Run `snitch` to launch the application
2. Set up AI backend (Ollama or Groq)
3. Set your current task
4. Start monitoring
5. Get back to work

## Build Commands

- `make build` - Build the binary
- `make install` - Install globally
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make build-all` - Build for all platforms

## License

MIT
