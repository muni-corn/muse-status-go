// Package sbattery is the "smart battery" pacakge, used for smart battery time
// estimates based on the user's weekly usage
package sbattery

import (
	"time"
)

// ChargeStatus acts as an enum for battery status
type ChargeStatus string

// Enum for battery statuses
const (
	Unknown     ChargeStatus = "Unknown"
	Discharging ChargeStatus = "Discharging"
	Charging    ChargeStatus = "Charging"
	Full        ChargeStatus = "Full"
)

const (
	recordInterval = 10 // record info every x minutes (status change overrides)
)

var (
	// battery icons
	dischargingIcons = [11]rune{'\uf08e', '\uf07a', '\uf07b', '\uf07c', '\uf07d', '\uf07e', '\uf07f', '\uf080', '\uf081', '\uf082', '\uf079'}
	chargingIcons    = [11]rune{'\uf89e', '\uf89b', '\uf086', '\uf087', '\uf088', '\uf89c', '\uf089', '\uf89d', '\uf08a', '\uf08b', '\uf085'}
	unknownIcon      = '\uf590'
	// nerd font icons
	// dischargingIcons = [11]rune{'\uf58d', '\uf579', '\uf57a', '\uf57b', '\uf57c', '\uf57d', '\uf57e', '\uf57f', '\uf580', '\uf581', '\uf578'}
	// chargingIcons    = [11]rune{'\uf585', '\uf586', '\uf587', '\uf588', '\uf589', '\uf58a', '\uf584'}
	// unknownIcon = '\uf590'
)

// returns a battery icon
func getBatteryIcon(status ChargeStatus, percentage int) rune {
	// get the battery icon
	var icon rune
	switch status {
	case Charging:
		chargingIndex := int((float32(percentage) / 100) * float32(len(chargingIcons)))
		if chargingIndex >= len(chargingIcons) {
			chargingIndex = len(chargingIcons) - 1
		}
		icon = chargingIcons[chargingIndex]
	case Discharging:
		dischargingIndex := int((float32(percentage) / 100) * float32(len(dischargingIcons)))
		if dischargingIndex >= len(dischargingIcons) {
			dischargingIndex = len(dischargingIcons) - 1
		}
		icon = dischargingIcons[dischargingIndex]
	case Full:
		// no display if full (return space character; found
		// that return the null character terminates i3bar's
		// json and will cause a problem
		return ' '
	}

	return icon
}

// formats a string telling at which time the battery will be empty/full
// e.g. "full at 3:30 pm"
func getTimeRemainingString(status ChargeStatus, timeDone time.Time) string {
	if (status == Charging || status == Discharging) && !timeDone.IsZero() {
		// get the time string prefix
		var timeStringPrefix string
		if status == Charging {
			timeStringPrefix = "full at "
		} else {
			timeStringPrefix = "until "
		}

		return timeStringPrefix + timeDone.Format("3:04 pm")
	}
	return ""
}

// SysPowerSupplyBaseDir is the base directory for power supply classes
const SysPowerSupplyBaseDir = "/sys/class/power_supply"

/*  DATA FILE FORMAT

data recorded like so:
key %/hour records

where "records" is the amount of times the parameter has been recorded.
used for recording a new average based on the current average and how
many times the parameter has been recorded before
for example:

C 3.14159 200

--- BEGIN FILE EXAMPLE ------------------------------------------------

C			// charging avg
C0			|
C1  		|
C2			| charging values by percentage
...			|
C9			|

D			// discharging avg
D0			|
D1			|
D2			| discharging avg values by hour of day
...			| (used for predicting nonexistent day-by-day values)
D23			|

S0			| sunday
S1			|
S2			| discharging values by hour by day of week
...			|
S23			|

M0			// monday
...

T0			// thursday
...

W0			// wednesday
...

R0			// thursday
...

F0			// friday
...

A0			// saturday
...
*/
