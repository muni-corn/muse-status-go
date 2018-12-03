package format

import (
	"strconv"
	"strings"
	"time"
)

// Chain chains status bites together, ensuring that there are no
// awkward spaces between bites.
func Chain(modules ...string) string {
	// start with some space at the beginning
	var final string

	for _, v := range modules {
		// trim space at the ends
		v = strings.TrimSpace(v)
		if v != "" {
			// add space between modules
			final += Separator() + v
		}
	}

	return final + Separator()
}

// Left aligns the original string to the left
func Left(original string) string {
	return "%{l}" + original
}

// Center aligns the original string to the center
func Center(original string) string {
	return "%{c}" + original
}

// Right aligns the original string to the right
func Right(original string) string {
	return "%{r}" + original
}

// Dim apples a half-transparent white color to the original string
func Dim(original string) string {
	return "%{F#80FFFFFF}" + original + "%{F-}"
}

// Warning renders the original string orange
func Warning(original string) string {
	return "%{F#FFAA00}" + original + "%{F-}"
}

// WarningBlink slowly blinks the original string oranage
func WarningBlink(original string) string {
	// convert unix nanoseconds to unix milliseconds
	milliseconds := time.Now().UnixNano() / 1000000

	// get alpha byte value
	x := float32((milliseconds/2)%1000) / 1000
	alpha := int(cubicEaseArc(x) * 255)

	// limit alpha minimum to 50%
	// we don't use math.Min because that would floop up the
	// animation
	alpha = alpha/2 + 255/2

	hex := ByteToHex(alpha)

	return "%{F#" + hex + "FFAA00}" + original + "%{F-}"
}

// FadeToDim colors the string according to interpolation from full white (0) to
// dim (1)
func FadeToDim(original string, interpolation float32) string {
	// constrain
	if interpolation < 0 {
		interpolation = 0
	} else if interpolation > 1 {
		interpolation = 1
	}

	// quintic graph
	x := interpolation * -1
	y := x*x*x*x*x + 1

	hex := ByteToHex(255/2 + int(y*255/2))
	return "%{F#" + hex + "FFFFFF}" + original + "%{F-}"
}

// Alert blinks the original string red
func Alert(original string) string {
	// convert nanoseconds to milliseconds
	milliseconds := time.Now().Nanosecond() / 1000000

	// get alpha byte value
	alpha := int(cubicEaseArc(float32(milliseconds)/1000) * 255)

	// limit alpha minimum to 25%
	alpha = alpha*3/4 + 255/4

	hex := ByteToHex(alpha)

	return "%{F#" + hex + "FF0000}" + original + "%{F-}"
}

// a cubic function except it's all concave down
func cubicEaseArc(x float32) float32 {
	x *= 2
	x--
	cubic := x * x * x
	if cubic > 0 {
		cubic *= -1
	}

	cubic++
	return cubic
}

// Separator returns 4 spaces as a separator between data
func Separator() string {
	return "    "
}

// ByteToHex takes a value from 0 to 255 and returns it in hexadecimal form
func ByteToHex(value int) string {
	// constrain to 0...255
	if value > 255 {
		return "FF"
	} else if value < 0 {
		return "00"
	}
	return strings.ToUpper(strconv.FormatInt(int64(value), 16))
}
