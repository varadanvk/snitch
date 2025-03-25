# Snitch: AI Productivity Watchdog

## Description

`snitch` is a local AI productivity watchdog that monitors your screen activity and calls you out when you're not focused. Using edge AI, it runs entirely on your device, analyzing your screen in real-time to determine if you're actually working. When it catches you slacking, it generates personalized notifications to get you back on track. Optional "snitch mode" lets you add accountability buddies who receive text updates about your procrastination habits.

## Key Features

- **Local AI Screen Monitoring**: Analyzes your screen content using vision-capable AI models without sending data to the cloud
- **Real-time Task Tracking**: Set your current task and get notifications when you go off-track
- **Sassy Personalized Notifications**: Receive custom messages to get you back on track
- **Accountability Buddy System**: Optional mode to let friends know when you're procrastinating
- **Privacy-Focused**: All processing done locally on your device
- **Modern UI**: Clean, minimal interface that stays out of your way
- **Productivity Analytics**: Track your focus patterns over time

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/snitch.git
   cd snitch
   ```

2. Create a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

4. Install Ollama (required) for AI-powered screen analysis:
   ```bash
   # Follow instructions at https://ollama.ai/
   # After installing, pull a vision-capable model:
   ollama pull llava
   # Start the Ollama service:
   ollama serve
   ```

5. (Optional) For accountability features, create a `.env` file with your Twilio credentials:
   ```
   TWILIO_ACCOUNT_SID=your_account_sid
   TWILIO_AUTH_TOKEN=your_auth_token
   TWILIO_PHONE_NUMBER=your_twilio_phone_number
   ```

## Usage

Run the application:
```bash
python src/main.py
```

### Basic Workflow

1. Enter your current task in the "What are you working on?" field
2. Click "Start Task" to begin monitoring
3. Get back to work!
4. Receive notifications when Snitch detects you're getting distracted
5. (Optional) If you've set up accountability buddies, they'll get text messages when you procrastinate too much

### Configuration

Snitch stores its configuration in `~/.snitch/config.json`, which includes:

- Productive/distracting app lists
- Notification intervals
- Accountability buddy information
- UI preferences

## Technical Details

### Architecture

Snitch uses a modular architecture with the following components:

- **Core**: Manages the application state and coordinates components
- **ML**: Handles screen analysis and activity classification
- **UI**: Provides the user interface
- **Utils**: Shared utility functions

### Technical Stack

- **Python**: Core language
- **CustomTkinter**: Modern UI toolkit
- **Ollama**: Local LLM for enhanced message generation
- **MSS**: Fast screen capture library
- **Twilio**: SMS integration for accountability features

## Project Structure

- `/src/main.py`: Main entry point, application initialization
- `/src/core/core.py`: Core functionality classes
- `/src/ml/ml.py`: ML and screen analysis components
- `/src/ui/ui.py`: User interface components
- `/src/utils/utils.py`: Utility functions and helpers
- `/requirements.txt`: Project dependencies
- `/assets/`: Application assets (icons, sounds, etc.)
- `/tests/`: Test suite

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request