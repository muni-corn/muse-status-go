package format

import (
	"fmt"
	"strings"
	"strconv"
)

var textFont, iconFont = "Roboto 10", "Material Design Icons 12"

// Mode is for different types of status modes, for different status bars that
// parse information differently
type Mode int

var mode Mode

// Definitions for Mode
const (
	LemonbarMode Mode = iota
	I3JSONMode
)

// Chain chains status bites together, ensuring that there are no
// awkward spaces between bites.
func Chain(blocks ...DataBlock) string {
	var first int
	var final string

	// huh. increment first until we find a module that
	// isn't nil or blank (empty for loop)
	for first = 0; first < len(blocks) && blocks[first] == nil; first++ { }

	// if everything is blank, return a blank string
	if first >= len(blocks) {
		return ""
	}

	switch mode {
	case I3JSONMode:
		final = I3JSONOf(blocks[first])
	default:
		final = LemonbarOf(blocks[first])
	}

	for i := first + 1; i < len(blocks); i++ {
		if blocks[i] == nil || blocks[i].Hidden() {
			continue
		}

		var v string
		switch mode {
		case I3JSONMode:
			v = I3JSONOf(blocks[i])
			v = Escape(v);
		default:
			v = LemonbarOf(blocks[i])
		}

		// trim space at the ends
		v = strings.TrimSpace(v)
		if v != "" {
			final += ModuleSeparator() + v
		}
	}

	return final
}

// Escape escapes characters for the i3 json protocol
func Escape(original string) string {
	return strings.ReplaceAll(original, "&", `&amp;`) // escape ampersand for json
}

// Action returns the original text wrapped in a clickable
// area
func Action(action, original string) string {
	if mode == I3JSONMode {
		return original
	}
	s := fmt.Sprintf("%%{A:%s:}%s%%{A}", action, original)
	return s
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

func interpolateColors(first, second Color, interpolation float32) (result Color, err error) {
	// rgbs
	firstInt, err := strconv.ParseInt(first.RGBHex, 16, 64)
	if err != nil {
		return
	}
	secondInt, err := strconv.ParseInt(second.RGBHex, 16, 64)
	if err != nil {
		return
	}
	r1, r2 := (firstInt>>16)&0xFF, (secondInt>>16)&0xFF
	g1, g2 := (firstInt>>8)&0xFF, (secondInt>>8)&0xFF
	b1, b2 := firstInt&0xFF, secondInt&0xFF

	// alphas
	a1, err := strconv.ParseInt(first.AlphaHex, 16, 64)
	if err != nil {
		return
	}
	a2, err := strconv.ParseInt(second.AlphaHex, 16, 64)
	if err != nil {
		return
	}

	a := int(float32(a1)*(1.0-interpolation) + float32(a2)*interpolation)
	r := int(float32(r1)*(1.0-interpolation) + float32(r2)*interpolation)
	g := int(float32(g1)*(1.0-interpolation) + float32(g2)*interpolation)
	b := int(float32(b1)*(1.0-interpolation) + float32(b2)*interpolation)

	resultRGBInt := r<<16 + g<<8 + b

	result.RGBHex = fmt.Sprintf("%06x", resultRGBInt)
	result.AlphaHex = fmt.Sprintf("%02x", a)
	return
}

// ModuleSeparator returns something that separates modules
// (spaces in Lemonbar mode, comma + space in i3 mode)
func ModuleSeparator() string {
	switch mode {
	case I3JSONMode:
		return ","
	case LemonbarMode:
		return "    "
	default:
		return "    "
	}
}

// Separator returns 4 spaces as a separator between data
func Separator() string {
	return "    "
}

// SetFormatMode sldksldksldkals;ala;skdkl;aslk
func SetFormatMode(m Mode) {
	mode = m
}

// GetFormatMode returns the format mode.
func GetFormatMode() Mode {
	return mode
}

// SetTextFont sets the regular text font
func SetTextFont(font string) {
	textFont = font
}

// SetIconFont sets the icon font
func SetIconFont(font string) {
	iconFont = font
}
