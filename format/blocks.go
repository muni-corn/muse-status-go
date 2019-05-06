package format

import (
	"fmt"
	"strings"
	"time"
)

// DataBlock is a block of data in the status bar
type DataBlock interface {
	Output() string
	Hidden() bool
}

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

// ClassicBlock has an icon, primary text, and dim secondary
// text.
type ClassicBlock struct {
	Name          string
	Icon          rune
	primaryText   string
	secondaryText string
	Urgency       Urgency
	hidden        bool
}

const (
	jsonTemplate      = `{"name":"%s","full_text":"%s","short_text":"%s","markup":"pango","separator":true}`
	pangoTemplate     = `<span color=\"#%s\" font=\"%s\">%s</span>`
	twoStringTemplate = "%s  %s"
)

// Output returns the ClassicBlock's output
func (c *ClassicBlock) Output() string {
	if c.hidden {
		return ""
	}

	defaultColor := primaryColor
	switch c.Urgency {
	case UrgencyLow:
		defaultColor = secondaryColor
	case UrgencyWarning:
		defaultColor = warningColor
	case UrgencyWarningPulse:
		// TODO
		defaultColor = warningColor
	case UrgencyAlarmPulse:
		defaultColor = getAlarmPulseColor()
	}

	var icon, primary, secondary string
	if c.Icon != '\x00' {
		icon = fmt.Sprintf(pangoTemplate, defaultColor, iconFont, string(c.Icon))
	}
	if strings.TrimSpace(c.primaryText) != "" {
		primary = fmt.Sprintf(pangoTemplate, defaultColor, textFont, c.primaryText)
	}
	if strings.TrimSpace(c.secondaryText) != "" {
		secondary = fmt.Sprintf(pangoTemplate, secondaryColor, textFont, c.secondaryText)
	}

	shortText := strings.TrimSpace(fmt.Sprintf(twoStringTemplate, icon, primary))

	var fullText string
	if c.secondaryText != "" {
		fullText = strings.TrimSpace(fmt.Sprintf(twoStringTemplate, shortText, secondary))
	} else {
		fullText = shortText
	}

	return fmt.Sprintf(jsonTemplate, c.Name, fullText, shortText)
}

// Hidden returns true if the ClassicBlock is hidden
func (c *ClassicBlock) Hidden() bool {
	return c.hidden
}

// SetHidden sets the ClassicBlock's visibility
func (c *ClassicBlock) SetHidden(h bool) {
	c.hidden = h
}

// Set sets the most common parameters of the ClassicBlock.
func (c *ClassicBlock) Set(urgency Urgency, icon rune, primaryText, secondaryText string) {
	c.Icon = icon
	c.Urgency = urgency
    c.SetPrimaryText(primaryText)
    c.SetSecondaryText(secondaryText)
}

// SetPrimaryText sets the primary text of the block
func (c *ClassicBlock) SetPrimaryText(text string) {
    c.primaryText = escape(text)
}

// SetSecondaryText sets the secondary text of the block
func (c *ClassicBlock) SetSecondaryText(text string) {
    c.secondaryText = escape(text)
}

func escape(text string) string {
    text = strings.ReplaceAll(text, "&", `&amp;`) // escape ampersand for json
    text = strings.ReplaceAll(text, "\"", `&quot;`) // escape quotes for json
    return text
}

func getAlarmPulseColor() (color string) {
	// convert nanoseconds to milliseconds
	milliseconds := time.Now().Nanosecond() / 1000000

	// get alpha byte value
    interpolation := cubicEaseArc(float32(milliseconds)/1000)
    
    color, err := interpolateColors(secondaryColor, alarmColor + "ff", interpolation)
    if err != nil {
        color = alarmColor
    }

    return 
}

// FadingBlock is a DataBlock with an icon and text that are
// dim by default. When its value changes, the data block
// lights up momentarily, then fades back to dim.
type FadingBlock struct {
	Name       string
	Icon       rune
	Text       string
	LastUpdate time.Time
	fading     bool
	hidden     bool
}

// Fading returns true if this FadingBlock is animating
func (f *FadingBlock) Fading() bool {
	return f.fading
}

// Set sets the common parameters of the ClassicBlock.
func (f *FadingBlock) Set(icon rune, text string) {
	f.Icon = icon
	f.Text = text
}

// Trigger starts the FadingBlock's animation
func (f *FadingBlock) Trigger() {
	f.fading = true
	f.LastUpdate = time.Now()
}

const secondsThreshold = 3

// Output returns the ClassicBlock's output
func (f *FadingBlock) Output() string {
	var color string
    // var err error
	if f.fading {
		secondsPassed := float32(time.Now().Sub(f.LastUpdate)) / float32(time.Second)
		x := secondsPassed / secondsThreshold
		x = x * x * x * x * x // quintic interpolation
        color, _ = interpolateColors(primaryColor, secondaryColor, x)
        // if err != nil {
        //     f.Text += "(err: " + err.Error() + ")"
        // }

		if secondsPassed > secondsThreshold {
			f.fading = false
		}
	} else {
		color = secondaryColor
	}

    color = strings.ToUpper(color)

	icon := fmt.Sprintf(pangoTemplate, color, iconFont, string(f.Icon))
	text := fmt.Sprintf(pangoTemplate, color, textFont, f.Text)
	// text := fmt.Sprintf(pangoTemplate, color, textFont, f.Text + " (#" + color + ")")

	shortText := icon
	fullText := fmt.Sprintf(twoStringTemplate, shortText, text)

	return fmt.Sprintf(jsonTemplate, f.Name, fullText, shortText)
}

// Hidden returns true if the ClassicBlock is hidden
func (f *FadingBlock) Hidden() bool {
	return f.hidden
}

// SetHidden sets the ClassicBlock's visibility
func (f *FadingBlock) SetHidden(h bool) {
	f.hidden = h
}
