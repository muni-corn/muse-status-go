package format

import (
	"time"
	"strconv"
	"strings"
)

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

// Alert blinks the original string red
func Alert(original string) string {
	// convert nanoseconds to milliseconds
	milliseconds := time.Now().Nanosecond() / 1000000

	// get alpha byte value
	alpha := int64(cubicEaseArc(float32(milliseconds) / 1000) * 255);

	// limit alpha to [64, 255]
	alpha = alpha * 3/4  + 255 / 4

	hex := strings.ToUpper(strconv.FormatInt(alpha, 16));

	return "%{F#" + hex + "FF0000}" + original + "%{F-}"
}

// a cubic function except it's all concave down
func cubicEaseArc(x float32) float32 {
	x *= 2
	x--
	cubic := x * x * x;
	if cubic > 0 {
		cubic *= -1
	}

	cubic++
	return cubic
}
