package scheduler

import (
	"app/handler"
	"time"
)

// StartAutomationScheduler starts the automation checker
func StartAutomationScheduler() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			handler.CheckAutomations()
		}
	}
}
