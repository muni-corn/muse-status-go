package date

import (
	"time"

	"github.com/muni-corn/muse-status/format"
)

// Block is a block that transmits time and date data
type Block struct {
	current    time.Time
	nextUpdate time.Time
}

// NewDateBlock returns a new date.Block
func NewDateBlock() *Block {
	b := &Block{}
	b.Update()
	return b
}

// NeedsUpdate returns true if the clock is out of date
func (b *Block) NeedsUpdate() bool {
	return b.current.After(b.nextUpdate)
}

// Update updates the clock
func (b *Block) Update() {
	b.current = time.Now()
	b.nextUpdate = b.current.Add(time.Minute).Truncate(time.Minute)
}

// Name returns "date"
func (b *Block) Name() string {
	return "date"
}

// Icon returns the clock icon
func (b *Block) Icon() rune {
	return icon
}

// Text returns the clock as primary, the date as secondary
func (b *Block) Text() (primary, secondary string) {
	return b.current.Format(timeFormat), b.current.Format(dateFormat)
}

// Fader returns nil; date module has no fader
func (b *Block) Fader() *format.Fader {
	return nil
}

// Hidden returns false; clock does not hide
func (b *Block) Hidden() bool {
	return false
}

// Urgency is always UrgencyNormal
func (b *Block) Urgency() format.Urgency {
	return format.UrgencyNormal
}

// ForceShort returns false; no force-shorting date yet
func (b *Block) ForceShort() bool {
	return false
}
