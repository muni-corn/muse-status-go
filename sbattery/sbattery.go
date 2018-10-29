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
)

// StartSmartBatteryBroadcast returns a channel that transfers intelligent
// battery information
func StartSmartBatteryBroadcast() chan string {
	// create a channel
	channel := make(chan string)

	// start the broadcast to it (async)
	go broadcast(channel)

	// return the channel
	return channel
}

func broadcast(channel chan string) string {
	for {
		channel <- status()
		time.Sleep(time.Second * 20)
	}
}

func status() string {
	output, err := exec.Command("acpi").Output()
	if err != nil {
		return "Error executing acpi. Is it installed?"
	}
	status, percentage, timeDone := parseReading(string(output))

	timeString := getTimeRemainingString(status, timeDone)

	finalOutput := getMainInfo(status, percentage) + timeString

	return finalOutput
}

// returns a colored icon and percentage, the main info of this
// module
func getMainInfo(status ChargeStatus, percentage int) string {
	// icone
	icon := getBatteryIcon(status, percentage)

	// base string
	base := icon + " " + strconv.Itoa(percentage) + "%  "

	switch status {
	case Charging:
		return base
	case Discharging:
		if percentage <= 15 {
			return format.Alert(base)
		} else if percentage <= 30 {
			return format.Warning(base)
		} else {
			return base
		}
	case Full:
		return icon + " Full"
	}

	// something's weird at this point
	return "Unknown"
}

// returns the battery icon
func getBatteryIcon(status ChargeStatus, percentage int) string {
	// get the battery icon
	var icon rune
	switch status {
	case Charging:
		// this is hell
		icon = chargingIcons[int((float32(percentage)/100)*float32(len(chargingIcons)))]
	case Discharging:
		icon = dischargingIcons[int((float32(percentage)/100)*float32(len(dischargingIcons)))]
	case Full:
		icon = chargingIcons[len(chargingIcons)-1]
	}

	return string(icon)
}

// returns a dimmed string telling at which time the battery will be empty/full
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

		return format.Dim(timeStringPrefix + timeDone.Format("3:04 pm"))
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

/*  DATA FILE FORMAT
 *
 *	data recorded in %/hour
 *	
 *	C			// charging avg
 *	C0			|
 *	C10  		|
 *	C20			| charging values by percentage
 *	...			| starting at value specified
 *	C90			|
 *	
 *	D			// discharging avg
 *	D0			|
 *	D1			|
 *	D2			| discharging avg values by hour of day
 *	...			| (used for predicting nonexistent day-by-day values)
 *	D23			|			
 *	
 *	S0			| sunday
 *	S1			|
 *	S2			| discharging values by hour by day of week
 *	...			|
 *	S23			|
 *	
 *	M0			// monday
 *	...
 *
 *	T0			// thursday
 *  ...
 *
 *	W0			// wednesday
 *  ...
 *
 *	R0			// thursday
 *  ...
 *
 *	F0			// friday
 *  ...
 *
 *	A0			// saturday 
 *  ...
 */
