package format

import (
	"strconv"
	"time"
)

var primaryColor = Color{ // {{{
	RGBHex: "ffffff",
	AlphaHex: "ff",
} // }}}
var secondaryColor = Color{ // {{{
	RGBHex: "ffffff",
	AlphaHex: "c0",
} // }}}
var transparentColor = Color { // {{{
	RGBHex: secondaryColor.RGBHex,
	AlphaHex: "00",
} // }}}
var warningColor = Color{ // {{{
	RGBHex: "ffaa00",
} // }}}
var alarmColor = Color{ // {{{
	RGBHex: "ff0000",
} // }}}

// PrimaryColor exposes the primaryColor to other packages
func PrimaryColor() Color {
	return primaryColor
}

// SecondaryColor exposes the secondaryColor to other packages
func SecondaryColor() Color {
	return secondaryColor
}

// Color represents a color in RRGGBB form. There is also an Alpha that can be
// provided separate from the rgb, as sometimes the alpha channel is either in
// argb or rgba
type Color struct {
	RGBHex string
	AlphaHex string
}

// HexString returns a hex representation of the string (no #), depending on
// the mode that muse-status is in
func (c *Color) HexString(mode Mode) string {
	switch(mode) {
	case LemonbarMode:
		return c.AlphaHex + c.RGBHex;
	default:
		return c.RGBHex + c.AlphaHex;
	}
}

// Colorer returns different colors for icon, primary, and
// secondary colors
type Colorer interface {
    IconColor() Color
    PrimaryColor() Color
    SecondaryColor() Color
}

// ByteToHex takes a value from 0 to 255 and returns it in hexadecimal form
func ByteToHex(value int) string {
	// constrain to 0...255
	if value > 255 {
		return "FF"
	} else if value < 0 {
		return "00"
	}
	return strconv.FormatInt(int64(value), 16)
}

// SetSecondaryColor sets the secondary (dim) color of
// muse-status.
func SetSecondaryColor(color string) {
	secondaryColor.RGBHex = color[:6];
	if len(color) == 8 {
		secondaryColor.AlphaHex = color[6:];
	} else if len(color) != 6 {
		println("invalid secondary color")
	}
	transparentColor = Color {
		RGBHex: secondaryColor.RGBHex,
		AlphaHex: "00",
	}
}

// SetPrimaryColor sets the primary color of
// muse-status.
func SetPrimaryColor(color string) {
	primaryColor.RGBHex = color[:6];
	if len(color) == 8 {
		primaryColor.AlphaHex = color[6:];
	} else if len(color) != 6 {
		println("invalid primary color")
	}
}

var (
	defCol = &defaultColorer{}
	wCol = &warnColorer{}
	aCol = &alarmColorer{}
	dimCol = &dimColorer{}

)

// GetDefaultColorer returns the default colorer
func GetDefaultColorer() Colorer {
	return defCol
}

// GetWarningColorer returns the pulsing warning colorer
func GetWarningColorer() Colorer {
	return wCol
}

// GetAlarmColorer returns the pulsing alarm colorer
func GetAlarmColorer() Colorer {
	return aCol
}

// GetDimColorer returns the pulsing dim colorer
func GetDimColorer() Colorer {
	return dimCol
}

// defaultColorer just returns the default colors {{{
type defaultColorer struct { }

// IconColor returns the default primaryColor
func (d defaultColorer) IconColor() Color {
	return primaryColor
}

// PrimaryColor returns the default primaryColor
func (d defaultColorer) PrimaryColor() Color {
	return primaryColor
}

// SecondaryColor returns the default secondaryColor
func (d defaultColorer) SecondaryColor() Color {
	return secondaryColor
}
// }}}

// dimColorer just returns the default secondaryColor for everything {{{
type dimColorer struct { }

// IconColor returns the default secondaryColor
func (d dimColorer) IconColor() Color {
	return secondaryColor
}

// PrimaryColor returns the default secondaryColor
func (d dimColorer) PrimaryColor() Color {
	return secondaryColor
}

// SecondaryColor returns the default secondaryColor
func (d dimColorer) SecondaryColor() Color {
	return secondaryColor
}
// }}}

// alarmColorer returns blinking red {{{
type alarmColorer struct { }

// IconColor returns blinking red
func (d alarmColorer) IconColor() Color {
	return getAlarmPulseColor()
}

// PrimaryColor returns blinking red
func (d alarmColorer) PrimaryColor() Color {
	return getAlarmPulseColor()
}

// SecondaryColor returns blinking red
func (d alarmColorer) SecondaryColor() Color {
	return getAlarmPulseColor()
}
// }}}

// warnColorer returns slow blinking orange {{{
type warnColorer struct { }

// IconColor returns slow blinking orange
func (d warnColorer) IconColor() Color {
	return getWarnPulseColor()
}

// PrimaryColor returns slow blinking orange
func (d warnColorer) PrimaryColor() Color {
	return getWarnPulseColor()
}

// SecondaryColor returns slow blinking orange
func (d warnColorer) SecondaryColor() Color {
	return getWarnPulseColor()
}
// }}}

// pulseColorer slowly blinks dim colors {{{
type pulseColorer struct { }

// IconColor returns slow blinking orange
func (d pulseColorer) IconColor() Color {
	return getDimPulseColor()
}

// PrimaryColor returns slow blinking orange
func (d pulseColorer) PrimaryColor() Color {
	return getDimPulseColor()
}

// SecondaryColor returns slow blinking orange
func (d pulseColorer) SecondaryColor() Color {
	return getDimPulseColor()
}
// }}}

func getAlarmPulseColor() Color { // {{{
	return getPulseColor(alarmColor, 1)
} // }}}

func getWarnPulseColor() Color { // {{{
	return getPulseColor(warningColor, 2)
} // }}}

func getDimPulseColor() Color { // {{{
	return getPulseColor(transparentColor, 3)
} // }}}

func getPulseColor(color Color, seconds float32) Color { // {{{
	var result Color

	if color.AlphaHex == "" {
		color.AlphaHex = "ff"
	}

	// get alpha byte value. interpolation is a value from
	// zero to one, calculated by unixMillis/maxMillis
	maxMillis := 1000 * seconds
	unixMillis := (time.Now().UnixNano() / int64(time.Millisecond)) % int64(maxMillis)
	interpolation := cubicEaseArc(float32(unixMillis) / maxMillis)

	result, err := interpolateColors(secondaryColor, color, interpolation)
	if err != nil {
		result = alarmColor
	}

	return result
} // }}}

// vim: foldmethod=marker
