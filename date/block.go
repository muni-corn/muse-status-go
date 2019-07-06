package date

import (
	"time"

	"github.com/muni-corn/muse-status/format"
)

// Block is a block that transmits time and date data
type Block struct {
	now    time.Time
	nextTime time.Time
	nextUpdate time.Time
}

// NewDateBlock returns a new date.Block
func NewDateBlock() *Block {
	b := &Block{}
	return b
}

// NextUpdateCheckTime is a function that does exactly what it says it does
func (b *Block) NextUpdateCheckTime() time.Time {
	return b.nextUpdate
}

// NeedsUpdate returns true if the clock is out of date
func (b *Block) NeedsUpdate() bool {
	b.nextTime = time.Now().Truncate(time.Minute)
	return b.now.Truncate(time.Minute) != b.nextTime
}

// Update updates the clock
func (b *Block) Update() {
	b.now = time.Now()
	b.nextUpdate = b.now.Add(time.Minute).Truncate(time.Minute)
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
	return b.now.Format(timeFormat), b.now.Format(dateFormat)
}

// Colorer returns the default colorer
func (b *Block) Colorer() format.Colorer {
	return format.GetDefaultColorer()
}

// Hidden returns false; clock does not hide
func (b *Block) Hidden() bool {
	return false
}

// ForceShort returns false; no force-shorting date yet
func (b *Block) ForceShort() bool {
	return false
}
