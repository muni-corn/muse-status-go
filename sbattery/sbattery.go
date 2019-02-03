// Package sbattery is the "smart battery" pacakge, used for smart battery time
// estimates based on the user's weekly usage
package sbattery

import (
	"muse-status/format"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ChargeStatus acts as an enum for battery status
type ChargeStatus int

// Enum for battery statuses
const (
	Unknown     ChargeStatus = 0
	Discharging ChargeStatus = 1
	Charging    ChargeStatus = 2
	Full        ChargeStatus = 3
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

// StartSmartBatteryBroadcast returns a channel that transfers intelligent
// battery information
func StartSmartBatteryBroadcast() chan *format.ClassicBlock {
	// create a channel
	channel := make(chan *format.ClassicBlock)
	block := &format.ClassicBlock{Name: "sbattery"}

	// start the broadcast to it (async)
	go broadcast(block, channel)

	// return the channel
	return channel
}

// an async function that broadcasts battery information to the specified
// channel
func broadcast(block *format.ClassicBlock, channel chan *format.ClassicBlock) {
	for {
		output, _ := exec.Command("acpi").Output()

		status, percentage, timeDone := parseReading(string(output))

		var urgency format.Urgency
		if status == Full {
			block.SetHidden(true)
			goto output
		} else if block.Hidden() {
			block.SetHidden(false)
		}

		switch {
		case percentage <= 15:
			urgency = format.UrgencyAlarmPulse
		case percentage <= 30:
			urgency = format.UrgencyWarning
		}

		block.Set(urgency, getBatteryIcon(status, percentage), strconv.Itoa(percentage)+"%", getTimeRemainingString(status, timeDone))

		output:
		channel <- block
		if (urgency == format.UrgencyAlarmPulse) {
			time.Sleep(time.Second / 15)
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

// returns the battery icon
func getBatteryIcon(status ChargeStatus, percentage int) rune {
	// get indices
	chargingIndex := int((float32(percentage) / 100) * float32(len(chargingIcons)))
	dischargingIndex := int((float32(percentage) / 100) * float32(len(dischargingIcons)))

	// constrain indices (but theoretically they should
	// never drop below zero)
	if chargingIndex >= len(chargingIcons) {
		chargingIndex = len(chargingIcons) - 1
	}
	if dischargingIndex >= len(dischargingIcons) {
		dischargingIndex = len(dischargingIcons) - 1
	}

	// get the battery icon
	var icon rune
	switch status {
	case Charging:
		icon = chargingIcons[chargingIndex]
	case Discharging:
		icon = dischargingIcons[dischargingIndex]
	case Full:
		// no display if full
		return '\000'
	}

	return icon
}

// returns a string telling at which time the battery will be empty/full
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

func parseReading(reading string) (status ChargeStatus, percentage int, timeDone time.Time) {
	// parse raw data
	split := strings.Split(reading, ", ")

	// get status
	rawStatus := split[0]
	switch {
	case strings.Contains(rawStatus, "Discharging"):
		status = Discharging
	case strings.Contains(rawStatus, "Charging"):
		status = Charging
	case strings.Contains(rawStatus, "Full"):
		status = Full
	default:
		status = Unknown
	}

	// get the percentage as an int
	if status == Full {
		percentage = 100
	} else {
		rawPercent := split[1]
		percentage, _ = strconv.Atoi(strings.TrimSuffix(rawPercent, "%"))
	}

	if len(split) >= 3 {
		rawTime := split[2]

		regex := regexp.MustCompile(`[^\d:]`)
		rawTime = regex.ReplaceAllString(rawTime, "")

		timeSplit := strings.Split(rawTime, ":")
		if len(timeSplit) == 3 {
			hours, _ := strconv.Atoi(timeSplit[0])
			minutes, _ := strconv.Atoi(timeSplit[1])
			seconds, _ := strconv.Atoi(timeSplit[2])

			// get the time at which the battery will be full/empty
			timeDone = time.Now().Add(time.Hour * time.Duration(hours)).Add(time.Minute * time.Duration(minutes)).Add(time.Second * time.Duration(seconds))
		}
	}

	return
}

// SysPowerSupplyBaseDir is the base directory for power supply
const SysPowerSupplyBaseDir = "/sys/class/power_supply"

func getBatteryStatus() ChargeStatus {
	return Unknown
}

func getBatteryPercentage() float32 {
	return 0
}

func takeReading() {

}

func recordReading() {

}

func getStringFromFile(filepath string) {

}

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
