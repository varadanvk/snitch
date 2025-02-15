## Snitch

### Description

`snitch` is a local AI productivity watchdog that monitors your screen activity and calls you out when you're not focused. Using edge AI, it runs entirely on your device, analyzing your screen in real-time to determine if you're actually working. When it catches you slacking, it generates sassy, personalized notifications to get you back on track. Optional "snitch mode" lets you add accountability buddies who receive text updates about your procrastination habits.

### Key Features:

- Local AI screen monitoring
- Real-time task tracking
- Accountability buddy system
- Privacy-focused (all processing done locally)
- System tray integration

### Technical Stack:

- Python desktop app
- Ollama for local LLM inference
- CustomTkinter for UI
- Screen capture and analysis
- Optional Twilio integration for accountability texts

Think of it as a judgmental AI assistant that lives on your desktop and isn't afraid to call you out for checking Twitter for the 15th time in an hour.

### Project Structure

`/src/main.py`

- Main entry point
- Initializes UI
- Sets up background processes
- Handles app lifecycle
- Error handling and logging setup

`/src/core/core.py`

- ScreenMonitor class: Handles screenshot capture and processing
- TaskManager class: Manages user goals and current tasks
- NotificationManager class: Handles local notifications
- AccountabilityManager class: Handles Twilio integration for snitching
- ConfigManager class: Manages app settings and preferences

`/src/ml/ml.py`

- OllamaInterface class: Handles communication with Ollama
- ScreenAnalyzer class: Processes screenshots to determine activity
- ActivityClassifier class: Determines if user is on/off task
- MessageGenerator class: Generates sassy notifications
- ActivityHistory class: Tracks and analyzes patterns of behavior

`/src/ui/ui.py`

- MainWindow class: Primary UI window using CustomTkinter
- TaskInputFrame class: UI for inputating tasks/goals
- SettingsFrame class: Configuration UI
- StatusFrame class: Shows current status and recent activities
- TrayIcon class: System tray integration
- BuddyManagerFrame class: UI for managing accountability buddies

`/src/utils/utils.py`

- Screenshot utilities (capture, save, process)
- Text processing helpers
- Time management utilities
- API wrappers for Twilio
- Config file handlers
- Custom exceptions
- Logging utilities

`/requirements.txt`

- customtkinter
- pillow
- numpy
- mss
- ollama
- requests
- twilio
- python-dotenv
- pynput

`/assets/`

- App icons
- UI elements
- Notification sounds
- Default avatars

`/tests/`

- test_core.py: Core functionality tests
- test_ml.py: ML component tests
- test_ui.py: UI component tests
- test_utils.py: Utility function tests
