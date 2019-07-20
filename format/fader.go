package format

import (
	"time"
)

// FadingColorer creates a fading effect between two colors.
type FadingColorer struct {
	Duration   float32
	StartColor Color
	EndColor   Color

	lastUpdate time.Time
	fading     bool
}

// Trigger starts the FadingBlock's animation
func (f *FadingColorer) Trigger() {
	f.fading = true
	f.lastUpdate = time.Now()
}

func (f *FadingColorer) IsFading() bool {
	return f.fading
}

func (f *FadingColorer) color() Color {
	var color Color

	if f.fading {
		secondsPassed := float32(time.Now().Sub(f.lastUpdate)) / float32(time.Second)
		x := secondsPassed / f.Duration
		x = x * x * x * x * x // quintic interpolation
		color, _ = interpolateColors(f.StartColor, f.EndColor, x)

		if secondsPassed > f.Duration {
			f.fading = false
		}
	} else {
		color = f.EndColor
	}

	return color
}

// IconColor returns the color of the fader
func (f *FadingColorer) IconColor() (color Color) {
	return f.color()
}

// PrimaryColor returns the color of the fader
func (f *FadingColorer) PrimaryColor() (color Color) {
	return f.color()
}

// SecondaryColor returns the color of the fader
func (f *FadingColorer) SecondaryColor() (color Color) {
	return f.color()
}

