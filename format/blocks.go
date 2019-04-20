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
    
    Fader() *Fader
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
    shortText := shortTextOf(b, primaryColor)

    // decide which fullText to use, in case we're forcing
    // short text
    var fullText string
    if b.ForceShort() {
        fullText = shortText
    } else {
        fullText = fullTextOf(b, primaryColor, secondaryColor)
    }

    // return the json
	return fmt.Sprintf(jsonTemplate, b.Name(), fullText, shortText)
}

func fullTextOf(b DataBlock, primaryColor, secondaryColor string) string {
	secondaryRaw, _ := b.Text()
    var secondary string
	if secondaryRaw != "" {
		secondary = fmt.Sprintf(pangoTemplate, secondaryColor, textFont, secondaryRaw)
	}

    return fmt.Sprintf(twoStringTemplate, shortTextOf(b, primaryColor), secondary)
}

func shortTextOf(b DataBlock, primaryColor string) string {
	iconRaw := b.Icon()
	primaryRaw, _ := b.Text()

    var icon, primary string

	if iconRaw != ' ' {
		icon = fmt.Sprintf(pangoTemplate, primaryColor, textFont, string(iconRaw))
	}
	if primaryRaw != "" {
		primary = fmt.Sprintf(pangoTemplate, primaryColor, textFont, strings.TrimSpace(primaryRaw))
	}

    return fmt.Sprintf(twoStringTemplate, icon, primary)
}

func getAlarmPulseColor() (color string) {
	// convert nanoseconds to milliseconds
	milliseconds := time.Now().Nanosecond() / 1000000

	// get alpha byte value
	interpolation := cubicEaseArc(float32(milliseconds) / 1000)

	color, err := interpolateColors(secondaryColor, alarmColor+"ff", interpolation)
	if err != nil {
		color = alarmColor
	}

	return
}
