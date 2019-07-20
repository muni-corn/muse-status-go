package brightness

func getIcon(percentage int) rune {
	index := percentage * len(brightnessIcons) / 100

	// constrain index (should never go below zero)
	if index >= len(brightnessIcons) {
		index = len(brightnessIcons) - 1
	}

	return brightnessIcons[index]
}
