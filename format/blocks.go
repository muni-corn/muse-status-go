package format

import (
	"fmt"
	"strings"
	"time"
)

// Urgency is a level of urgency applied to a block
type Urgency int

// Urgency definitions
const (
	UrgencyNormal Urgency = iota
	UrgencyLow
	UrgencyWarning
	UrgencyWarningPulse
	UrgencyAlarmPulse
)

// DataBlock is a piece of data in the status bar.
type DataBlock interface {
	NeedsUpdate() bool
	Update()

	Name() string
	Icon() rune
	Text() (primary, secondary string)

	Colorer() Colorer
	Hidden() bool
	Urgency() Urgency
	ForceShort() bool
}

const (
	jsonTemplate      = `{"name":"%s","full_text":"%s","short_text":"%s","markup":"pango","separator":true}`
	pangoTemplate     = `<span color=\"#%s\" font=\"%s\">%s</span>`
	twoStringTemplate = "%s  %s"
)

// JSONOf Block b. Turns the information of b into a JSON
// object for the i3 status protocol
func JSONOf(b DataBlock) string {
	if b.Hidden() {
		return ""
	}

	// get short text
	shortText := shortTextOf(b)

	// decide which fullText to use, in case we're forcing
	// short text
	var fullText string
	if b.ForceShort() {
		fullText = shortText
	} else {
		fullText = fullTextOf(b)
	}

	// return the json
	return fmt.Sprintf(jsonTemplate, b.Name(), fullText, shortText)
}

func fullTextOf(b DataBlock) string {
	secondaryRawText, _ := b.Text()
	var secondaryText string
	if secondaryRawText != "" {
		secondaryText = fmt.Sprintf(pangoTemplate, b.Colorer().SecondaryColor(), textFont, secondaryRawText)
	}

	return fmt.Sprintf(twoStringTemplate, shortTextOf(b), secondaryText)
}

func shortTextOf(b DataBlock) string {
	iconRaw := b.Icon()
	primaryRawText, _ := b.Text()

	var icon, primaryText string

	if iconRaw != ' ' {
		icon = fmt.Sprintf(pangoTemplate, b.Colorer().IconColor(), textFont, string(iconRaw))
	}
	if primaryRawText != "" {
		primaryText = fmt.Sprintf(pangoTemplate, b.Colorer().PrimaryColor(), textFont, strings.TrimSpace(primaryRawText))
	}

	return fmt.Sprintf(twoStringTemplate, icon, primaryText)
}

func getAlarmPulseColor() string {
	return getPulseColor(alarmColor, 1)
}

func getWarnPulseColor() string {
	return getPulseColor(warningColor, 2)
}

func getPulseColor(color string, seconds float32) string {
	var result string

	// color must be 8 characters long to be valid. if there
	// are only 6, we'll append "ff" for our 100% opacity
	// alpha value.
	if len(color) == 6 {
		color += "ff"
	} else if len(color) != 8 {
		return color
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
}
