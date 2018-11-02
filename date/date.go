package date

import (
	f "muse-status/format"
	"time"
)

const (
	timeFormat = "3:04 pm"
	dateFormat = "Mon, Jan 2"
)

// StartDateBroadcast creates a string channel that transmits the current date
func StartDateBroadcast() chan string {
	channel := make(chan string)

	go func() {
		var lastTimeString string

		for {
			// get current time
			now := time.Now()
			timeString := now.Format("3:04 pm")

			// default sleep interval to a 20th of a second
			sleepInterval := time.Second / 20

			if timeString != lastTimeString {
				dateString := now.Format("Mon, Jan 2")
				
				// output to channel
				channel <- "\uf150  " + timeString + "  " + f.Dim(dateString)

				// if the time has changed, we're not changing again anytime
				// soon. get number of seconds until next minute change and
				// sleep for that many seconds minus 0.25 seconds (in case we
				// sleep too long; this allows for an accurate time change as
				// we're updating every 20th second when we're anticipating a
				// time change)
				sleepInterval = time.Second * time.Duration(60 - now.Second()) - time.Second / 4

				// update lastTimeString
				lastTimeString = timeString

				// skip sleeping if we're not going to sleep anyways
				if sleepInterval < 0 {
					continue
				}
			}

			// sleep
			time.Sleep(sleepInterval)
		}
	}()

	return channel
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
