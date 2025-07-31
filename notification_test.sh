#!/bin/bash

echo "Testing macOS notifications..."

# Test 1: Direct osascript notification
echo "Test 1: Direct macOS notification"
osascript -e 'display notification "This is a test notification" with title "Snitch Test" sound name "Ping"'

# Test 2: Wait a moment
sleep 2

# Test 3: Another notification
echo "Test 2: Another notification"
osascript -e 'display notification "If you see this, notifications are working!" with title "Success" sound name "Ping"'

echo "Tests completed. Check if you saw notifications appear on your screen." 