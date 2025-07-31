# Snitch - AI Productivity Watchdog

## Project Structure
```
src/
├── main.py          # Application entry point
├── core/            # Core functionality (SnitchCore, managers)
├── ui/              # User interface (ModernUI)
├── ml/              # Machine learning (ScreenAnalyzer, ActivityClassifier)
├── utils/           # Utility functions
├── templates/       # HTML templates
└── static/          # CSS/JS assets
```

## Build/Lint/Test Commands
```bash
# Run the application
python src/main.py

# Run all tests
python -m pytest tests/

# Run a single test file
python -m pytest tests/test_core.py

# Run a specific test class
python -m pytest tests/test_core.py::TestSnitchCore

# Run a specific test method
python -m pytest tests/test_core.py::TestSnitchCore::test_initialization
```

## Code Style Guidelines

### Imports
- Standard library imports first
- Third-party imports next
- Local imports last
- Organize imports in logical groups

### Naming Conventions
- Class names: PascalCase (e.g., `SnitchCore`)
- Functions/methods: snake_case (e.g., `capture_screen`)
- Constants: UPPER_SNAKE_CASE (e.g., `DEFAULT_CONFIG_DIR`)
- Private members: prefixed with underscore (e.g., `_load_config`)

### Formatting & Documentation
- Use 4 spaces for indentation (no tabs)
- Comprehensive docstrings for classes and functions
- Google-style docstrings with Args/Returns sections
- Inline comments for complex logic
- Maximum line length: 88 characters (Black default)

### Type Hinting
- Use type hints for all function parameters and return values
- Import types from the `typing` module when needed
- Use `Optional`, `List`, `Dict`, etc. appropriately

### Error Handling
- Use try/except blocks around critical operations
- Create custom exception classes for specific error types
- Log errors with appropriate context and log levels
- Avoid generic except clauses

### Class Design
- Follow single responsibility principle
- Maintain clear separation of concerns
- Use dependency injection for callbacks and external dependencies
- Keep classes focused and cohesive

## Development Setup
1. Create virtual environment: `python -m venv venv`
2. Activate virtual environment: `source venv/bin/activate` (Unix) or `venv\Scripts\activate` (Windows)
3. Install dependencies: `pip install -r requirements.txt`
4. Install Ollama and pull llava model: `ollama pull llava`

## Notes
- All processing runs locally for privacy
- Configuration stored in `~/.snitch/config.json`
- Uses CustomTkinter for modern UI
- Integrates with Twilio for SMS notifications (optional)