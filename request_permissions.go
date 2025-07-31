package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🔔 Snitch Notification Permission Setup")
	fmt.Println("=====================================")

	// Step 1: Request notification permissions
	fmt.Println("\n1. Requesting notification permissions...")

	// Use osascript to request permissions by sending a test notification
	cmd := exec.Command("osascript", "-e", `
		display dialog "Snitch needs notification permissions to alert you about productivity. Click OK to test notifications." buttons {"OK", "Cancel"} default button "OK"
	`)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Could not request permissions:", err)
		return
	}

	// Check if user clicked OK
	if string(output) == "OK" {
		fmt.Println("✅ User granted permission")
	} else {
		fmt.Println("❌ User denied permission")
		return
	}

	// Step 2: Test notification
	fmt.Println("\n2. Testing notification...")
	time.Sleep(1 * time.Second)

	testCmd := exec.Command("osascript", "-e",
		`display notification "If you see this notification, permissions are working correctly!" with title "Snitch Test" sound name "Ping"`)

	if err := testCmd.Run(); err != nil {
		fmt.Println("❌ Test notification failed:", err)
		return
	}

	fmt.Println("✅ Test notification sent!")

	// Step 3: Send a few more test notifications
	fmt.Println("\n3. Sending additional test notifications...")

	notifications := []string{
		"Productive activity detected - great work!",
		"Distraction detected - consider refocusing",
		"Task reminder - how is your progress?",
	}

	for i, msg := range notifications {
		time.Sleep(2 * time.Second)
		fmt.Printf("   Sending notification %d/3...\n", i+1)

		cmd := exec.Command("osascript", "-e",
			fmt.Sprintf(`display notification "%s" with title "Snitch Test" sound name "Ping"`, msg))

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Notification %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Notification %d sent!\n", i+1)
		}
	}

	fmt.Println("\n🎉 Setup complete!")
	fmt.Println("If you saw the notifications, Snitch should work properly.")
	fmt.Println("If you didn't see notifications, check System Preferences > Notifications & Focus")
}
