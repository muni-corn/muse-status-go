package date

import (
	"time"
)

const (
	timeFormat = "3:04 pm"
	dateFormat = "Mon, Jan 2"
	// icon = '\uf64f' // nerd font icon
	icon = '\uf150'
)

// GetGreeting returns a greeting based on the hour of the day
func GetGreeting() string {
	hour := time.Now().Hour()

	switch {
	case hour < 12:
		return "Good morning!"
	case hour >= 12 && hour < 17:
		return "Good afternoon!"
	default:
		return "Good evening!"
	}
}
