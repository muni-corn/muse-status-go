// Package sbattery is the "smart battery" pacakge, used for
// smart battery time estimates based on the user's hourly
// usage
package sbattery

// TODO split these functions into more files. this is too
// much to be in one file.

import (
	"errors"
	"io/ioutil"
	"muse-status/format"
	"muse-status/utils"
	"strconv"
	"strings"
	"time"
)

type batteryReading struct {
	time       time.Time
	percentage float32
}

type rateData struct {
	minutesPerPercent float32
	numReadings       int
}

// ChargeStatus acts as an enum for battery status
type chargeStatus int

// Enum for battery statuses
const (
	Unknown chargeStatus = iota
	Discharging
	Charging
	Full
)

const (
	recordInterval = 5   // record info every x minutes (status change overrides)
	numRecordLimit = 500 // maximum number of records to keep for each key
	baseDir        = "/sys/class/power_supply/"
	battery        = "BAT0"
)

var (
	// battery icons
	dischargingIcons = [11]rune{'\uf08e', '\uf07a', '\uf07b', '\uf07c', '\uf07d', '\uf07e', '\uf07f', '\uf080', '\uf081', '\uf082', '\uf079'}
	chargingIcons    = [11]rune{'\uf89e', '\uf89b', '\uf086', '\uf087', '\uf088', '\uf89c', '\uf089', '\uf89d', '\uf08a', '\uf08b', '\uf085'}

	// rate information
	hourlyDischargingAverage   rateData
	hourlyChargingAverage      rateData
	chargingAvgByPercent       [10]rateData
	dischargingAvgByHourOfDay  [24]rateData
	dischargingAvgByHourOfWeek [7 * 24]rateData
)

// StartSmartBatteryBroadcast returns a channel that transfers intelligent
// battery information
func StartSmartBatteryBroadcast() chan string {
	// create a channel
	channel := make(chan string)
	go func() {
		var lastStatus chargeStatus
		var lastReading batteryReading
		var lastRecordedReading batteryReading

		info, alert := "", false

		for {
			currentStatus, err := getBatteryStatus()
			currentReading := takeReading()

			// record the delta if we're past due for a
			// recording
			shouldRecordReading := currentReading.time.Sub(lastRecordedReading.time) >= time.Minute*recordInterval
			if shouldRecordReading {
				recordReading(currentStatus, lastRecordedReading, currentReading)
			}

			if err != nil {
				channel <- format.Dim("Couldn't get battery info")
			} else if lastReading != currentReading || lastStatus != currentStatus || alert {
				// update status and reading
				lastStatus = currentStatus
				lastReading = currentReading

				info, alert = getFullInfo(currentReading, currentStatus)
				channel <- info

				// if battery is low, sleep for less time
				// and continue the loop for animation
				// (we continue so we don't sleep for 3
				//  seconds)
				if alert {
					time.Sleep(time.Second / 15)
					continue
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()

	// return the channel
	return channel
}

// returns the colored icon, percentage, and time remaining
func getFullInfo(reading batteryReading, status chargeStatus) (string, bool) {
	timeDone, err := getDueTime(status)

	// get time remaining as string. if there was an error
	// in calculating the time until the battery is
	// full/empty, we
	var timeString string
	if err != nil {
		timeString = ""
	} else {
		timeString = getTimeString(status, timeDone)
	}

	// get the final, full information string
	finalOutput := getMainInfo(status, reading.percentage) + timeString

	return finalOutput, reading.percentage <= 15
}

// returns a colored icon and percentage, the main info of this
// module
func getMainInfo(status chargeStatus, percentage float32) string {
	// icon
	icon := getBatteryIcon(status, percentage)

	// base string
	base := icon + " " + strconv.Itoa(int(percentage)) + "%  "

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
		return "\uf084"
	}

	// something's weird at this point
	return "\uf091"
}

// returns an appropriate battery icon as a string based on
// battery status and percentage
func getBatteryIcon(status chargeStatus, percentage float32) string {
	// get icon indices
	chargingIndex := int((percentage / 100) * float32(len(chargingIcons)))
	dischargingIndex := int((percentage / 100) * float32(len(dischargingIcons)))

	// constrain indices (but they should never drop below
	// zero)
	if chargingIndex >= len(chargingIcons) {
		chargingIndex = len(chargingIcons) - 1
	}
	if dischargingIndex >= len(dischargingIcons) {
		dischargingIndex = len(dischargingIcons) - 1
	}

	// gets the icon from the index in the correct icon
	// array
	var icon rune
	switch status {
	case Charging:
		icon = chargingIcons[chargingIndex]
	case Discharging:
		icon = dischargingIcons[dischargingIndex]
	case Full:
		icon = chargingIcons[len(chargingIcons)-1]
	}

	return string(icon)
}

// returns a dimmed string telling at which time the battery
// will be empty/full, e.g. "full at 3:30 pm".
func getTimeString(status chargeStatus, timeDone time.Time) string {
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

// returns a time object, where the time is when the battery
// will be either empty or full
func getDueTime(status chargeStatus) (time.Time, error) {
	switch status {
	case Charging:
		return getChargingDueTime()
	case Discharging:
		return getDischargingDueTime()
	}
	return time.Time{}, errors.New("Invalid status")
}

func getChargingDueTime() (time.Time, error) {
	var durationLeft time.Duration

	// here we set up a simulated battery percentage that
	// we'll use to calculate a time estimate
	percProjection, err := getBatteryPercentage()
	if err != nil {
		return time.Time{}, err
	}

	// while our simulated percentage is less than 100
	for percProjection < 100 {
		// get index that will get our rate information at
		// this percentage point
		index := int(percProjection / 10)
		rateAtThisPerc := chargingAvgByPercent[index].minutesPerPercent

		// this is the amount of percentage until the next
		// multiple of ten
		percToNextTen := float32(index*10+1) - percProjection

		// if there is no rate information for this
		// particular multiple of ten, just use the overall
		// average. if that doesn't exist, we'll return a
		// not-enough-information error
		if rateAtThisPerc != 0 {
			durationLeft += time.Duration(percToNextTen * rateAtThisPerc * float32(time.Minute))
			// TODO DRY this out
		} else if hourlyChargingAverage.minutesPerPercent != 0 {
			durationLeft += time.Duration(percToNextTen * hourlyChargingAverage.minutesPerPercent * float32(time.Minute))
		} else {
			return time.Time{}, errors.New("Not enough information to get time remaining")
		}
	}

	// adjust for over-compensation
	// NOTE this may not be necessary but whatever lol
	if percProjection > 100 {
		var rate float32
		lastChargingRate := chargingAvgByPercent[len(chargingAvgByPercent)-1].minutesPerPercent
		if lastChargingRate != 0 {
			rate = lastChargingRate
			// TODO DRY this out:
		} else if hourlyChargingAverage.minutesPerPercent != 0 {
			rate = hourlyChargingAverage.minutesPerPercent
		} else {
			return time.Time{}, errors.New("Not enough information to get time remaining")
		}
		durationLeft -= time.Duration((percProjection - 100) * rate * float32(time.Minute))
	}

	return time.Now().Add(durationLeft), nil
}

func getDischargingDueTime() (time.Time, error) {
	var durationLeft time.Duration

	// here we set up a simulated battery percentage that
	// we'll use to calculate a time estimate
	percProjection, err := getBatteryPercentage()
	simulatedTime := time.Now()
	if err != nil {
		return time.Time{}, err
	}

	// while our simulated percentage is greater than 0
	for percProjection > 0 {
		nextHour := simulatedTime.Add(time.Hour).Truncate(time.Hour)

		// get index that will get our rate information at
		// this hour of the week and day
		weekIndex := getHourOfWeek(simulatedTime)
		rateThisHourOfWeek := dischargingAvgByHourOfWeek[weekIndex].minutesPerPercent
		dayIndex := getHourOfDay(simulatedTime)
		rateThisHourOfDay := dischargingAvgByHourOfWeek[dayIndex].minutesPerPercent

		// minutes until the next hour
		minutesToNextHour := nextHour.Sub(simulatedTime) / time.Minute
		durationLeft += time.Duration(minutesToNextHour) * time.Minute

		// if there is no rate information for this
		// particular multiple of ten, just use the overall
		// average. if that doesn't exist, we'll return a
		// not-enough-information error
		// TODO DRY this out
		if rateThisHourOfWeek != 0 {
			percProjection -= float32(minutesToNextHour) / rateThisHourOfWeek
		} else if rateThisHourOfDay != 0 {
			percProjection -= float32(minutesToNextHour) / rateThisHourOfDay
		} else if hourlyDischargingAverage.minutesPerPercent != 0 {
			percProjection -= float32(minutesToNextHour) / hourlyDischargingAverage.minutesPerPercent
		}

		// advance time
		simulatedTime = nextHour
	}

	// adjust for over-compensation
	if percProjection < 0 {
		// soooo the time at which the percentage dropped
		// below zero was during the previous hour that
		// simulatedTime is at this point. how many minutes
		// ago was the percentage at zero?

		// get rate to use, using fallbacks as needed
		var rate float32

		timeOneHourAgo := simulatedTime.Add(-time.Hour)

		weekIndex := getHourOfWeek(timeOneHourAgo)
		dayIndex := getHourOfDay(timeOneHourAgo)

		weekRate := dischargingAvgByHourOfWeek[weekIndex].minutesPerPercent
		dayRate := dischargingAvgByHourOfDay[dayIndex].minutesPerPercent
		hourlyRate := hourlyDischargingAverage.minutesPerPercent
		if weekRate != 0 {
			rate = weekRate
		} else if dayRate != 0 {
			rate = dayRate
		} else if hourlyRate != 0 {
			rate = hourlyRate
		} else {
			return time.Time{}, errors.New("Not enough information to calculate depletion time")
		}

		// basically = deltaPerc * rate
		minutesAgo := percProjection * rate
		durationLeft -= time.Duration(minutesAgo) * time.Minute
	}

	return time.Now().Add(durationLeft), nil
}

func getBatteryStatus() (status chargeStatus, err error) {
	// Unknown by default
	status = Unknown

	str, err := getStringFromFile(baseDir + battery + "/status")
	if err == nil {
		// parse the status
		switch {
		case strings.Contains(str, "Discharging"):
			status = Discharging
		case strings.Contains(str, "Charging"):
			status = Charging
		case strings.Contains(str, "Full"):
			status = Full
		}
	}
	return
}

func getBatteryPercentage() (value float32, err error) {
	chargeNow, err := utils.GetIntFromFile(baseDir + battery + "/charge_now")

	if err == nil {
		chargeFull, err := utils.GetIntFromFile(baseDir + battery + "/charge_full")
		if err == nil {
			value = float32(chargeNow) / float32(chargeFull)
		}
	}

	return
}

func takeReading() batteryReading {
	percentage, err := getBatteryPercentage()
	if err != nil {
		return batteryReading{}
	}
	return batteryReading{time: time.Now(), percentage: percentage}
}

func recordReading(status chargeStatus, last batteryReading, current batteryReading) {
	if status != Charging && status != Discharging {
		return
	}

	minutesPassed := current.time.Sub(last.time) / time.Minute
	percentagePassed := current.percentage - last.percentage
	rate := float32(minutesPassed) / percentagePassed

	switch {
	case rate >= 0:
		numReadingsBefore := float32(hourlyChargingAverage.numReadings)
		avgBefore := float32(hourlyChargingAverage.minutesPerPercent)
		newAvg := avgBefore*numReadingsBefore/(numReadingsBefore+1) + rate/(numReadingsBefore+1)
		hourlyChargingAverage.minutesPerPercent = newAvg
		hourlyChargingAverage.numReadings++

		startTen := int(last.percentage / 10)
		endTen := int(current.percentage / 10)

		for i := startTen; i != (endTen+1)%10; i = (i + 1) % 10 {
			numReadingsBefore := float32(hourlyDischargingAverage.numReadings)
			avgBefore := float32(hourlyDischargingAverage.minutesPerPercent)

			newAvg := avgBefore*numReadingsBefore/(numReadingsBefore+1) + rate/(numReadingsBefore+1)

			hourlyDischargingAverage.minutesPerPercent = newAvg
			hourlyDischargingAverage.numReadings++
		}

	case rate < 0:
		/* TODO modularize this crap */

		// update overall average (hda = hourlyDischargingAverage) //
		numReadingsBefore := float32(hourlyDischargingAverage.numReadings)
		avgBefore := float32(hourlyDischargingAverage.minutesPerPercent)
		newAvg := avgBefore*numReadingsBefore/(numReadingsBefore+1) + rate/(numReadingsBefore+1)
		hourlyDischargingAverage.minutesPerPercent = newAvg
		hourlyDischargingAverage.numReadings++
		/////////////////////////////////////////////////////////////

		// update by hour of day //
		startHourOfDay := getHourOfDay(last.time)
		endHourOfDay := getHourOfDay(current.time)
		for i := startHourOfDay; i != (endHourOfDay+1)%10; i = (i + 1) % 10 {
			numReadingsBefore := float32(dischargingAvgByHourOfDay[i].numReadings)
			avgBefore := float32(dischargingAvgByHourOfDay[i].minutesPerPercent)
			newAvg := avgBefore*numReadingsBefore/(numReadingsBefore+1) + rate/(numReadingsBefore+1)
			dischargingAvgByHourOfDay[i].minutesPerPercent = newAvg
			dischargingAvgByHourOfDay[i].numReadings++
		}
		///////////////////////////

		// update by hour of week //
		startHourOfWeek := getHourOfWeek(last.time)
		endHourOfWeek := getHourOfWeek(current.time)
		for i := startHourOfWeek; i != (endHourOfWeek+1)%7*24; i = (i + 1) % 7 * 24 {
			numReadingsBefore := float32(dischargingAvgByHourOfWeek[i].numReadings)
			avgBefore := float32(dischargingAvgByHourOfWeek[i].minutesPerPercent)
			newAvg := avgBefore*numReadingsBefore/(numReadingsBefore+1) + rate/(numReadingsBefore+1)
			dischargingAvgByHourOfWeek[i].minutesPerPercent = newAvg
			dischargingAvgByHourOfWeek[i].numReadings++
		}
		////////////////////////////
	}
}

func getHourOfDay(t time.Time) int {
	return t.Hour()
}

func getHourOfWeek(t time.Time) int {
	return int(t.Weekday())*24 + getHourOfDay(t)
}

// TODO
func saveData() {

}

func getStringFromFile(filepath string) (str string, err error) {
	output, err := ioutil.ReadFile(filepath)
	if err == nil {
		str = strings.TrimSpace(string(output))
	}
	return
}

/*  DATA FILE FORMAT

data recorded like so:
key minutes/% records

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

W0			|
W1			|
W2			| discharging values by hour of week
...			|
W167		|

*/
