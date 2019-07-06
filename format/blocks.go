package format

import (
	"fmt"
	"strings"
	"time"
)

// // Urgency is a level of urgency applied to a block
// type Urgency int

// // Urgency definitions
// const (
// 	UrgencyNormal Urgency = iota
// 	UrgencyLow
// 	UrgencyWarning
// 	UrgencyWarningPulse
// 	UrgencyAlarmPulse
// )

// DataBlock is a piece of data in the status bar.
type DataBlock interface {
	NextUpdateCheckTime() time.Time
	NeedsUpdate() bool
	Update()

	Name() string
	Icon() rune
	Text() (primary, secondary string)

	Colorer() Colorer
	Hidden() bool
	// Urgency() Urgency
	ForceShort() bool
}

const (
	jsonTemplate      = `{"name":"%s","full_text":"%s","short_text":"%s","markup":"pango","separator":true}`
	pangoTemplate     = `<span color=\"#%s\" font=\"%s\">%s</span>`
	twoStringTemplate = "%s  %s"
	lemonTemplate     = "%{%s} %s"
)

// Alignment specifies where on lemonbar a block should be aligned.
type Alignment string

// Alignment definitions
const (
	Right Alignment = "r"
	Center = "c"
	Left = "l"
)

// LemonbarOf a block. returns a string representation of the block that can be
// parsed by lemonbar
func LemonbarOf(b DataBlock, align Alignment) string {
	if b.Hidden() {
		return ""
	}

	primary, secondary := b.Text()
	icon := string(b.Icon())

	// color first
	c := b.Colorer()
	if c != nil {
		pColor := c.PrimaryColor()
		sColor := c.SecondaryColor()
		iColor := c.IconColor()
		icon = fmt.Sprintf(lemonTemplate, "F#"+iColor.AlphaHex+iColor.RGBHex, icon)
		primary = fmt.Sprintf(lemonTemplate, "F#"+pColor.AlphaHex+pColor.RGBHex, primary)
		secondary = fmt.Sprintf(lemonTemplate, "F#"+sColor.AlphaHex+sColor.RGBHex, secondary)
	}

	// then align
	return fmt.Sprintf("%{%s} %s", string(align), icon + "  " + primary + "  " + secondary + "  ");
	
}

// I3JSONOf Block b. Turns the information of b into a JSON
// object for the i3 status protocol
func I3JSONOf(b DataBlock) string {
	if b.Hidden() {
		return ""
	}

	// get short text
	shortText := shortPangoOf(b)

	// decide which fullText to use, in case we're forcing
	// short text
	var fullText string
	if b.ForceShort() {
		fullText = shortText
	} else {
		fullText = fullPangoOf(b)
	}

	// return the json
	return fmt.Sprintf(jsonTemplate, b.Name(), fullText, shortText)
}

func fullPangoOf(b DataBlock) string {
	secondaryRawText, _ := b.Text()
	var secondaryText string
	if secondaryRawText != "" {
		secondaryText = fmt.Sprintf(pangoTemplate, b.Colorer().SecondaryColor(), textFont, secondaryRawText)
	}

	return fmt.Sprintf(twoStringTemplate, shortPangoOf(b), secondaryText)
}

func shortPangoOf(b DataBlock) string {
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

func getAlarmPulseColor() Color {
	return getPulseColor(alarmColor, 1)
}

func getWarnPulseColor() Color {
	return getPulseColor(warningColor, 2)
}

func getDimPulseColor() Color {
	return getPulseColor(transparentColor, 3);
}

func getPulseColor(color Color, seconds float32) Color {
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
}
