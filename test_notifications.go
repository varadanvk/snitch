package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🔔 Snitch Notification Test")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Println("This will test if notifications work on your system.")
	fmt.Println("If you don't see notifications, you may need to:")
	fmt.Println("1. Go to System Preferences > Notifications & Focus")
	fmt.Println("2. Find your terminal app (Terminal, iTerm, VS Code)")
	fmt.Println("3. Enable notifications for it")
	fmt.Println()

	// Send a test notification
	fmt.Println("Sending test notification...")
	cmd := exec.Command("osascript", "-e",
		`display notification "This is a test notification from Snitch. If you see this, notifications are working!" with title "Snitch Test" sound name "Ping"`)

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Notification failed:", err)
		fmt.Println("This might be a permissions issue.")
		return
	}

	fmt.Println("✅ Test notification sent!")
	fmt.Println()

	// Send a few more notifications with delays
	notifications := []string{
		"Productive activity detected - great work!",
		"Distraction detected - consider refocusing on your task",
		"Task reminder - how is your progress?",
	}

	for i, msg := range notifications {
		time.Sleep(3 * time.Second)
		fmt.Printf("Sending notification %d/3...\n", i+1)

		cmd := exec.Command("osascript", "-e",
			fmt.Sprintf(`display notification "%s" with title "Snitch Test" sound name "Ping"`, msg))

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Notification %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Notification %d sent!\n", i+1)
		}
	}

	fmt.Println()
	fmt.Println("🎉 Test complete!")
	fmt.Println("If you saw the notifications, Snitch should work properly.")
	fmt.Println("If you didn't see notifications, check your notification settings.")
}
