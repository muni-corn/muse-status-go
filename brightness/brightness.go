package brightness

import (
	"github.com/muni-corn/muse-status/utils"
	"strconv"
	"time"
)

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
