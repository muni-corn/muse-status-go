package volume

import (
	"github.com/muni-corn/muse-status/format"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var (
	volumeIcons     = [3]rune{'', '', ''}
	muteIcon        = ''
	percentageRegex = regexp.MustCompile(`\[(\d*?)%?\]`)  // matches the volume percentage as an int
	onOffRegex      = regexp.MustCompile(`\[([A-z]*?)\]`) // matches 'on' or 'off' within square brackets
)

// StartVolumeBroadcast returns a channel that is fed audio volume information
func StartVolumeBroadcast() chan *format.FadingBlock {
	channel := make(chan *format.FadingBlock)
	block := &format.FadingBlock{Name: "volume"}

	go func() {
		lastVolume := -2

		// loop
		for {
			// get current volume
			current, err := getCurrentVolume()

			// check for an error, continue if there is one
			if err != nil {
				println("Error getting volume: " + err.Error())
				time.Sleep(2 * time.Second)
				continue
			} 

			// if the volume has changed
			if current != lastVolume {

				var icon rune
				var text string
				if current <= 0 {
					icon = muteIcon
					text = "Muted"
				} else {
					icon = getIcon(current)
					text = strconv.Itoa(current) + "%"
				}

				block.Set(icon, text)
				block.Trigger()
				channel <- block

				// update old data
				lastVolume = current
			} 

			// animate
			if block.Fading() {
				// faster framerate
				channel <- block
				time.Sleep(time.Second / 20)
			} else {
				time.Sleep(time.Second / 5)
			}
		}
	}()

	return channel
}

// returns the current volume percentage as an int, or zero
// if muted
func getCurrentVolume() (percentage int, err error) {
	output, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		return
	}

	strOutput := string(output)

	mixerStatus := onOffRegex.FindStringSubmatch(strOutput)[1] // should be 'on' or 'off'. if it's not, then wtf
	if mixerStatus == "off" {
		percentage = 0 // muted
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
