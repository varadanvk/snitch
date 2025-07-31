//quick script to test if notifications are working

package main

import (
	"fmt"
	"log"
	"time"

	"snitch-tui/src/core"
	"snitch-tui/src/notifications"
)

func main() {
	fmt.Println("Starting notification test...")
	
	manager := notifications.NewManager(time.Duration(1) * time.Second)
	
	fmt.Println("Sending test notification...")
	err := manager.SendActivityNotification(core.Activity{
		Activity:     "Test Activity",
		Application:  "Test App",
		IsProductive: true,
	})
	
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	} else {
		fmt.Println("Notification sent successfully!")
	}
	
	// Also test custom notification
	fmt.Println("Sending custom notification...")
	err = manager.SendCustomNotification("Test Title", "This is a test notification")
	if err != nil {
		log.Printf("Error sending custom notification: %v", err)
	} else {
		fmt.Println("Custom notification sent successfully!")
	}
}
