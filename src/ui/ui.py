import customtkinter as ctk
from typing import Callable, Optional
from datetime import datetime

# Configure the default theme for a Notion-like appearance
NOTION_COLORS = {
    "bg_light": "#ffffff",
    "bg_dark": "#191919",
    "sidebar_light": "#fbfbfa",
    "sidebar_dark": "#202020",
    "text_primary_light": "#37352f",
    "text_primary_dark": "#ffffff",
    "text_secondary_light": "#787774",
    "text_secondary_dark": "#999999",
    "accent": "#2eaadc",
    "hover_light": "#f1f1ef",
    "hover_dark": "#2d2d2d",
}


class ModernFrame(ctk.CTkFrame):
    """A modern frame with hover effects and better styling."""

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.configure(
            corner_radius=8,
            border_width=1,
            border_color=("gray85", "gray25"),
            fg_color=("gray95", "gray10"),
        )


class SidebarButton(ctk.CTkButton):
    """A modern sidebar button with hover effects."""

    def __init__(self, *args, **kwargs):
        kwargs.update(
            {
                "corner_radius": 8,
                "height": 40,
                "border_spacing": 10,
                "fg_color": "transparent",
                "text_color": ("gray25", "gray85"),
                "hover_color": ("gray90", "gray20"),
                "anchor": "w",
            }
        )
        super().__init__(*args, **kwargs)


class MainWindow(ctk.CTkFrame):
    def __init__(self, master, task_callback: Callable[[str], None], **kwargs):
        super().__init__(master, **kwargs)
        self.task_callback = task_callback

        # Configure the main window
        self.configure(fg_color=("gray95", "gray10"))

        # Create the layout
        self.create_sidebar()
        self.create_main_content()

    def create_sidebar(self):
        """Create the sidebar with navigation and settings."""
        self.sidebar = ctk.CTkFrame(
            self,
            width=250,
            corner_radius=0,
            fg_color=("gray98", "gray15"),
        )
        self.sidebar.pack(side="left", fill="y", padx=0, pady=0)
        self.sidebar.pack_propagate(False)

        # Logo and title
        self.logo_frame = ctk.CTkFrame(self.sidebar, fg_color="transparent")
        self.logo_frame.pack(fill="x", padx=15, pady=(15, 5))

        self.title = ctk.CTkLabel(
            self.logo_frame,
            text="Snitch",
            font=ctk.CTkFont("Inter", size=24, weight="bold"),
            text_color=("gray15", "gray90"),
        )
        self.title.pack(side="left", padx=5)

        # Navigation buttons
        self.nav_frame = ctk.CTkFrame(self.sidebar, fg_color="transparent")
        self.nav_frame.pack(fill="x", padx=10, pady=10)

        nav_items = [
            ("🎯 Tasks", self.show_tasks),
            ("📊 Analytics", self.show_analytics),
            ("⚡ Focus Mode", self.toggle_focus_mode),
            ("⚙️ Settings", self.show_settings),
        ]

        for text, command in nav_items:
            btn = SidebarButton(
                self.nav_frame,
                text=text,
                command=command,
                font=ctk.CTkFont("Inter", size=14),
            )
            btn.pack(fill="x", padx=5, pady=2)

    def create_main_content(self):
        """Create the main content area."""
        self.main_content = ctk.CTkFrame(self, fg_color="transparent")
        self.main_content.pack(side="left", fill="both", expand=True, padx=20, pady=20)

        # Task input area
        self.task_frame = ModernFrame(self.main_content)
        self.task_frame.pack(fill="x", pady=(0, 20))

        self.task_label = ctk.CTkLabel(
            self.task_frame,
            text="What are you working on?",
            font=ctk.CTkFont("Inter", size=18, weight="bold"),
            text_color=("gray15", "gray90"),
        )
        self.task_label.pack(padx=20, pady=(20, 10), anchor="w")

        self.task_entry = ctk.CTkEntry(
            self.task_frame,
            placeholder_text="Enter your current task...",
            height=45,
            font=ctk.CTkFont("Inter", size=14),
            border_width=1,
            corner_radius=8,
        )
        self.task_entry.pack(fill="x", padx=20, pady=(0, 10))

        self.task_button = ctk.CTkButton(
            self.task_frame,
            text="Start Task",
            font=ctk.CTkFont("Inter", size=14, weight="bold"),
            height=40,
            corner_radius=8,
            command=self.set_task,
        )
        self.task_button.pack(padx=20, pady=(0, 20))

        # Status area
        self.status_frame = ModernFrame(self.main_content)
        self.status_frame.pack(fill="x", pady=(0, 20))

        self.status_label = ctk.CTkLabel(
            self.status_frame,
            text="Ready to start monitoring...",
            font=ctk.CTkFont("Inter", size=14),
            height=60,
        )
        self.status_label.pack(padx=20, pady=20)

        # Recent activities
        self.activities_frame = ModernFrame(self.main_content)
        self.activities_frame.pack(fill="both", expand=True)

        self.activities_label = ctk.CTkLabel(
            self.activities_frame,
            text="Recent Activities",
            font=ctk.CTkFont("Inter", size=16, weight="bold"),
            text_color=("gray15", "gray90"),
        )
        self.activities_label.pack(padx=20, pady=(20, 10), anchor="w")

        # Placeholder for activities list
        self.activities_list = ctk.CTkTextbox(
            self.activities_frame,
            font=ctk.CTkFont("Inter", size=13),
            wrap="word",
            height=200,
        )
        self.activities_list.pack(fill="both", expand=True, padx=20, pady=(0, 20))

    def set_task(self):
        """Handle setting a new task."""
        task = self.task_entry.get()
        if task:
            self.task_callback(task)
            self.task_entry.delete(0, "end")
            self.update_status(f"🎯 Currently working on: {task}")
            self.add_activity(f"Started task: {task}")

    def update_status(self, text: str):
        """Update the status display."""
        self.status_label.configure(
            text=text, font=ctk.CTkFont("Inter", size=14, weight="bold")
        )

    def add_activity(self, text: str):
        """Add an activity to the activities list."""
        timestamp = datetime.now().strftime("%H:%M")
        self.activities_list.insert("1.0", f"{timestamp} - {text}\n")

    # Navigation callbacks
    def show_tasks(self):
        """Show the tasks view."""
        self.update_status("Viewing tasks...")

    def show_analytics(self):
        """Show the analytics view."""
        self.update_status("Viewing analytics...")

    def toggle_focus_mode(self):
        """Toggle focus mode."""
        self.update_status("Focus mode toggled...")

    def show_settings(self):
        """Show the settings view."""
        self.update_status("Viewing settings...")


class ModernUI:
    def __init__(self, root: ctk.CTk):
        self.root = root
        self.setup_window()
        self.main_window = None

    def setup_window(self):
        """Configure the main window."""
        self.root.title("Snitch")
        self.root.geometry("1200x800")

        # Configure the default appearance
        ctk.set_appearance_mode("dark")
        ctk.set_default_color_theme("blue")

        # Configure window style
        self.root.configure(fg_color=("gray95", "gray10"))

    def initialize(self, task_callback: Callable[[str], None]):
        """Initialize the UI with the main window."""
        self.main_window = MainWindow(self.root, task_callback)
        self.main_window.pack(fill="both", expand=True)

    def update_status(self, text: str):
        """Update the status display."""
        if self.main_window:
            self.main_window.update_status(text)

    def add_activity(self, text: str):
        """Add an activity to the log."""
        if self.main_window:
            self.main_window.add_activity(text)
