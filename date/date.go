package date

import (
	f "muse-status/format"
	"time"
)

// StartDateBroadcast takes a string channel and inputs to it every second
func StartDateBroadcast(channel chan string) {
	go func() {
		for {
			// get current time
			now := time.Now()
			channel <- "\uf150  " + now.Format("3:04 pm") + "  " + f.Dim(now.Format("Mon, Jan 2"))

			// sleep for a second
			time.Sleep(time.Second)
		}
	}()
}

// GetGreeting returns a greeting based on the hour of the day
func GetGreeting() string {
	hour := time.Now().Hour()

	switch {
	case hour >= 4 && hour < 12:
		return "Good morning!"
	case hour >= 12 && hour < 17:
		return "Good afternoon!"
	default:
		return "Good evening!"
	}
}
