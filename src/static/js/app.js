// Initialize Socket.IO
const socket = io();

// DOM Elements
const taskInput = document.getElementById("task-input");
const statusDisplay = document.getElementById("status-display");
const activityLog = document.getElementById("activity-log");

// Handle setting a new task
function setTask() {
  const task = taskInput.value.trim();
  if (task) {
    socket.emit("set_task", { task });
    taskInput.value = "";
  }
}

// Handle task input with Enter key
taskInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    setTask();
  }
});

// Socket event handlers
socket.on("status_update", (data) => {
  statusDisplay.textContent = data.message;
});

socket.on("activity_update", (data) => {
  const activityItem = document.createElement("div");
  activityItem.className = "py-2 px-3 bg-gray-50 dark:bg-gray-700 rounded-lg";

  const timestamp = document.createElement("span");
  timestamp.className = "text-sm text-gray-500 dark:text-gray-400";
  timestamp.textContent = new Date().toLocaleTimeString() + " - ";

  const message = document.createElement("span");
  message.className = "text-gray-700 dark:text-gray-300";
  message.textContent = data.message;

  activityItem.appendChild(timestamp);
  activityItem.appendChild(message);

  activityLog.insertBefore(activityItem, activityLog.firstChild);
});

// Theme handling
function setTheme(isDark) {
  if (isDark) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
}

// Check system theme preference
if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
  setTheme(true);
}

// Listen for system theme changes
window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", (e) => {
    setTheme(e.matches);
  });
