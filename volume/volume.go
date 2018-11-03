package volume

import (
	"muse-status/format"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var (
	volumeIcons     = [3]rune{'', '', ''}
	muteIcon        = ''
	percentageRegex = regexp.MustCompile(`\[(\d*?)%?\]`)      // matches the volume percentage as an int
	onOffRegex      = regexp.MustCompile(`\[([A-z]*?)\]`) // matches 'on' or 'off' within square brackets
)

// StartVolumeBroadcast returns a channel that is fed audio volume information
func StartVolumeBroadcast() chan string {
	channel := make(chan string)

	go func() {
		lastVolume := 0
		lastChangeTime := int64(0)
		interpolation := float32(0.0)

		// loop
		for {
			// get current volume percentage
			current, err := getCurrentVolume()

			// check for an error
			if err != nil {
				channel <- format.Dim("Error getting volume: " + err.Error())
			} else {
				// if the brightness has changed
				if current != lastVolume {
					// current time
					now := time.Now().UnixNano()

					// update out-of-date data
					lastChangeTime = now
					lastVolume = current
					interpolation = 0
				}

				// animate
				if interpolation < 1 {

					// TODO getInterpolation(): in format? this is also used in
					// the brightness package. could be included in the
					// FadeToDim function
					now := time.Now().UnixNano()
					interpolation = float32(now-lastChangeTime) / float32(int(time.Second)*2)

					var status string
					if current < 0 {
						status = string(muteIcon) + "  Muted"
					} else {
						status = string(getIcon(current)) + "  " + strconv.Itoa(current) + "%"
					}
					channel <- format.FadeToDim(status, interpolation)
					time.Sleep(time.Second / 15)
				} else {
					time.Sleep(time.Second / 2)
				}
			}
		}
	}()

	return channel
}

// returns the current volume percentage as an int, or returns a negative number
// if muted
func getCurrentVolume() (percentage int, err error) {
	output, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		percentage = -2
		return
	}

	strOutput := string(output)

	mixerStatus := onOffRegex.FindStringSubmatch(strOutput)[1] // should be 'on' or 'off'. if it's not, then wtf
	if mixerStatus == "off" {
		percentage = -1
	} else if mixerStatus == "on" {
		percentageStr := percentageRegex.FindStringSubmatch(strOutput)[1]
		percentage, err = strconv.Atoi(percentageStr)
	}
	return
}

func getIcon(percentage int) rune {
	index := percentage * len(volumeIcons) / 100

	// constrain index (should never go below zero)
	if index >= len(volumeIcons) {
		index = len(volumeIcons) - 1
	}

	return volumeIcons[index]
}
