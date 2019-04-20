package format

import (
    "time"
)

// Fader creates a fading effect between two colors.
type Fader struct {
    Duration float32
    StartColor string
    EndColor string

    lastUpdate time.Time
    fading bool
}

// Trigger starts the FadingBlock's animation
func (f *Fader) Trigger() {
	f.fading = true
	f.lastUpdate = time.Now()
}

// Output returns the color of the fader
func (f *Fader) Output() (color string) {
	if f.fading {
		secondsPassed := float32(time.Now().Sub(f.lastUpdate)) / float32(time.Second)
		x := secondsPassed / f.Duration
		x = x * x * x * x * x // quintic interpolation
		color, _ = interpolateColors(primaryColor, secondaryColor, x)

		if secondsPassed > f.Duration {
			f.fading = false
		}
	} else {
		color = secondaryColor
	}

	// color = strings.ToUpper(color)

	return color
}
