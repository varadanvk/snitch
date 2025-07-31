#!/bin/bash

# Snitch Installation Script
set -e

echo "Installing Snitch..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.19+ first."
    exit 1
fi

# Install the binary
go install github.com/varadanvk/snitch@latest

echo "✓ Snitch installed successfully!"
echo ""
echo "To get started:"
echo "1. Run 'snitch' to launch the application"
echo "2. Set up your AI backend (Ollama or Groq)"
echo "3. Set your current task and start monitoring"
echo ""
echo "For Ollama setup: ollama pull llava"
echo "For Groq setup: Get API key from https://console.groq.com/"