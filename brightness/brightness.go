package brightness

import (
	"github.com/muni-corn/muse-status/format"
	"strconv"
	"time"
	"github.com/muni-corn/muse-status/utils"
)

const (
	baseDir = "/sys/class/backlight/"
)

var (
	brightnessIcons = [6]rune{'', '', '', '', '', ''}
	// brightnessIcons = [6]rune{'\uf5da', '\uf5db', '\uf5dc', '\uf5dd', '\uf5de', '\uf5df'} // nerd font icons
	card            = "amdgpu_bl0"
)

// StartBrightnessBroadcast returns a string channel that is fed screen
// brightness information
func StartBrightnessBroadcast() chan *format.FadingBlock {
	channel := make(chan *format.FadingBlock)
	block := &format.FadingBlock{Name: "brightness"}

	go func() {
		max, err := getMaxBrightness()
		lastBrightness := 0

		if err != nil {
			println("Brightness encountered a fatal initialization error: ", err.Error())
			return
		}

		// loop
		for {
			// get current brightness (in whatever the heck
			// units Linux uses, not a percentage)
			current, err := getCurrentBrightness()

			// check for an error, continue if there is one
			if err != nil {
				println("Brightness encountered an error: ", err.Error())
				time.Sleep(2 * time.Second)
				continue
			} 

			// get the brightness percentage from 0 to 100
			brightnessPercentage := current * 100 / max

			// if the brightness has changed, update things
			if brightnessPercentage != lastBrightness {
				icon := getIcon(brightnessPercentage)
				text := strconv.Itoa(brightnessPercentage) + "%"

				block.Set(icon, text)
				block.Trigger()
				channel <- block

				// update old data
				lastBrightness = brightnessPercentage
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
