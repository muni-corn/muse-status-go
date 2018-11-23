package brightness

import (
	"io/ioutil"
	"muse-status/format"
	"strconv"
	"strings"
	"time"
	"muse-status/utils"
)

const (
	baseDir = "/sys/class/backlight/"
)

var (
	brightnessIcons = [6]rune{'', '', '', '', '', ''}
	card            = "amdgpu_bl0"
)

// StartBrightnessBroadcast returns a string channel that is fed screen
// brightness information
func StartBrightnessBroadcast() chan string {
	channel := make(chan string)

	go func() {
		max, err := getMaxBrightness()
		lastBrightness := 0
		lastChangeTime := int64(0)
		interpolation := float32(0.0)

		// if there's an error, fohgetaboutit
		if err != nil {
			channel <- ""
		}

		// loop
		for {
			// get current brightness
			current, err := getCurrentBrightness()

			// check for an error
			if err != nil {
				channel <- format.Dim("Error getting brightness")
			} else {
				// get the brightness percentage from 0 to 100
				brightnessPercentage := current * 100 / max

				// if the brightness has changed
				if brightnessPercentage != lastBrightness {
					// current time
					now := time.Now().UnixNano()

					// update out-of-date data
					lastChangeTime = now
					lastBrightness = brightnessPercentage
					interpolation = 0
				}

				// animate
				if interpolation < 1 {
					now := time.Now().UnixNano()
					interpolation = float32(now-lastChangeTime) / float32(int(time.Second)*2)
					status := string(getIcon(brightnessPercentage)) + "  " + strconv.Itoa(brightnessPercentage) + "%"
					channel <- format.FadeToDim(status, interpolation)
					time.Sleep(time.Second / 20)
				} else {
					time.Sleep(time.Second / 5)
				}
			}
		}
	}()

	return channel
}

func getMaxBrightness() (value int, err error) {
	return utils.GetIntFromFile(baseDir + card + "/max_brightness")
}

func getCurrentBrightness() (value int, err error) {
	return utils.GetIntFromFile(baseDir + card + "/brightness")
}

func getIcon(percentage int) rune {
	index := percentage * len(brightnessIcons) / 100

	// constrain index (should never go below zero)
	if index >= len(brightnessIcons) {
		index = len(brightnessIcons) - 1
	}

	return brightnessIcons[index]
}
