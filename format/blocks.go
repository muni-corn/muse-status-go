package format

import (
	"fmt"
	"strings"
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

const (
	pangoTemplate     = `<span color="#%s" font="%s">%s</span>`
	twoStringTemplate = "%s  %s"
	lemonTemplate     = "%%{%s} %s"
)

// DataBlock is a piece of data in the status bar.
type DataBlock interface {
	StartBroadcast() <-chan bool // returns a channel that sends signals to update the status bar
	Update()

	Name() string
	Hidden() bool
	ForceShort() bool

	Output(mode Mode) string
}

type ClassicBlock interface {
	DataBlock

	Text() (primary, secondary string)
	Icon() rune
	Colorer() Colorer
}

// BanneringBlock has the ability to display banners in the status bar
type BanneringBlock interface {
	Banner(interpolation float32) string
	Name() string // used to update banners
	Activate()
}

// LemonbarOf a block. returns a string representation of the block that can be
// parsed by lemonbar
func LemonbarOf(b ClassicBlock) string {
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

		// println("block: " + b.Name())
		// println("icon: " + icon)
		// println("primary: " + primary)
		// println("secondary: " + secondary)
		icon = fmt.Sprintf(lemonTemplate, "F#"+iColor.AlphaHex+iColor.RGBHex, icon)
		primary = fmt.Sprintf(lemonTemplate, "F#"+pColor.AlphaHex+pColor.RGBHex, primary)
		secondary = fmt.Sprintf(lemonTemplate, "F#"+sColor.AlphaHex+sColor.RGBHex, secondary)
	}

	// then align
	return icon + "  " + primary + "  " + secondary
}

type I3JSONBlock struct {
	Name      string `json:"name"`
	FullText  string `json:"full_text"`
	ShortText string `json:"short_text"`
	Markup    string `json:"markup"`
	Separator bool   `json:"separator"`
}

// I3JSONOf Block b. Turns the information of b into a JSON
// object for the i3 status protocol
func I3JSONOf(b ClassicBlock) *I3JSONBlock {
	if b.Hidden() {
		return nil
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

	j := I3JSONBlock{
		Name:      b.Name(),
		FullText:  fullText,
		ShortText: shortText,
		Markup:    "pango",
		Separator: true,
	}

	return &j
}

func fullPangoOf(b ClassicBlock) string {
	_, secondaryRawText := b.Text()
	var secondaryText string
	if secondaryRawText != "" {
		secondaryText = fmt.Sprintf(pangoTemplate, b.Colorer().SecondaryColor().HexString(mode), textFont, secondaryRawText)
	}

	return fmt.Sprintf(twoStringTemplate, shortPangoOf(b), secondaryText)
}

func shortPangoOf(b ClassicBlock) string {
	iconRaw := b.Icon()
	primaryRawText, _ := b.Text()

	var icon, primaryText string

	if iconRaw != ' ' {
		icon = fmt.Sprintf(pangoTemplate, b.Colorer().IconColor().HexString(mode), iconFont, string(iconRaw))
	}
	if primaryRawText != "" {
		primaryText = fmt.Sprintf(pangoTemplate, b.Colorer().PrimaryColor().HexString(mode), textFont, strings.TrimSpace(primaryRawText))
	}

	return fmt.Sprintf(twoStringTemplate, icon, primaryText)
}

// vim: foldmethod=syntax
