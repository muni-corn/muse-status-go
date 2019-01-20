package format

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var secondaryColor = "FFFFFFFF"

// Chain chains status bites together, ensuring that there are no
// awkward spaces between bites.
func Chain(modules ...string) string {
	var final string
	firstNonBlank := true

	for _, v := range modules {
		// trim space at the ends
		v = strings.TrimSpace(v)
		if v != "" {
			if firstNonBlank {
				final += v
				firstNonBlank = false
			} else {
				final += Separator() + v
			}
		}
	}

	return final + Separator()
}

// Action returns the original text wrapped in a clickable
// area
func Action(action, original string) string {
	s := fmt.Sprintf("%%{A:%s:}%s%%{A}", action, original)
	return s
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
	return "%{F#C0" + secondaryColor + "}" + original + "%{F-}"
}

// Warning renders the original string orange
func Warning(original string) string {
	return "%{F#FFAA00}" + original + "%{F-}"
}

// WarningBlink slowly blinks the original string oranage
func WarningBlink(original string) string {
	if original == "" {
		return ""
	}

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
	if original == "" {
		return ""
	}

	// constrain
	if interpolation < 0 {
		interpolation = 0
	} else if interpolation > 1 {
		interpolation = 1
	}

	// quintic graph. apparently this is backwards
	x := interpolation * -1
	y := x*x*x*x*x + 1

	color, err := interpolateColors("C0"+ secondaryColor, "FFFFFFFF", y)
	if err != nil {
		println(err.Error())
		color = secondaryColor
	}

	return "%{F#" + color + "}" + original + "%{F-}"
}

// Alert blinks the original string red
func Alert(original string) string {
	if original == "" {
		return ""
	}

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

func interpolateColors(first, second string, interpolation float32) (result string, err error) {
	firstInt, err := strconv.ParseInt(first, 16, 64)
	if err != nil {
		return
	}

	secondInt, err := strconv.ParseInt(second, 16, 64)
	if err != nil {
		return
	}

	r1, r2 := firstInt&0xFF, secondInt&0xFF
	g1, g2 := (firstInt>>8)&0xFF, (secondInt>>8)&0xFF
	b1, b2 := (firstInt>>16)&0xFF, (secondInt>>16)&0xFF
	a1, a2 := (firstInt>>24)&0xFF, (secondInt>>24)&0xFF

	r := int(float32(r1)*(1-interpolation) + float32(r2)*interpolation)
	g := int(float32(g1)*(1-interpolation) + float32(g2)*interpolation)
	b := int(float32(b1)*(1-interpolation) + float32(b2)*interpolation)
	a := int(float32(a1)*(1-interpolation) + float32(a2)*interpolation)

	resultInt := a<<24 + b<<16 + g<<8 + r

	result = strconv.FormatInt(int64(resultInt), 16)
	return
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

// SetSecondaryColor sets the secondary (dim) color of
// muse-status.
func SetSecondaryColor(color string) {
	switch {
	case len(color) == 6:
		secondaryColor = color
	default:
		println("Invalid secondary color. Defaulting to white.")
	}
}
