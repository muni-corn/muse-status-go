package brightness

import (
    "fmt"

    "github.com/muni-corn/muse-status/format"
)

const (
	baseDir = "/sys/class/backlight/"
)

var (
	brightnessIcons = [6]rune{'', '', '', '', '', ''}
	// brightnessIcons = [6]rune{'\uf5da', '\uf5db', '\uf5dc', '\uf5dd', '\uf5de', '\uf5df'} // nerd font icons
	card            = "amdgpu_bl0"
)

var maxBrightness int

func init() {
    maxBrightness, err := getMaxBrightness()
    if err != nil {
        panic(err)
    }
}

// Block is a block that contains device
// brightness information
type Block struct {
    format.DataBlock

    Card string
    currentBrightness int
    lastBrightness int

    text string
    icon rune
    fader *format.Fader
}

// NewBrightnessBlock returns a new brightness.Block
func NewBrightnessBlock(card string) *Block {
    return &Block{
        Card: card, 
        fader: &format.Fader{Duration: 3},
    }
}

// Name returns the name "brightness"
func (b *Block) Name() string {
    return "brightness"
}

// NeedsUpdate returns true if the Block needs its text
// updated
func (b *Block) NeedsUpdate() bool {
    var err error
    b.currentBrightness, err = getCurrentBrightness()

    if err == nil && b.currentBrightness != b.lastBrightness {
        b.lastBrightness = b.currentBrightness
        return true
    }

    return false
}

// Update updates the text and icon of the block
func (b *Block) Update() {
    b.text = fmt.Sprintf("%d%", b.currentBrightness * 100 / maxBrightness)
    b.icon = getIcon(b.currentBrightness)
    b.fader.Trigger()
}

// Icon returns the brightness icon
func (b *Block) Icon() rune {
    return b.icon
}

// Text returns the text of the block
func (b *Block) Text() (primary, secondary string) {
    return b.text, ""
}

// Fader returns a pointer to the block's fader, for color
func (b *Block) Fader() *format.Fader {
    return b.fader
}

// Hidden always returns false because brightness never has
// a need to be hidden
func (b *Block) Hidden() bool {
    return false
}

// Urgency always returns Urgency, since brightness module's
// color is controlled by its fader
func (b *Block) Urgency() format.Urgency {
    return format.UrgencyNormal
}

// ForceShort returns false because we never really need to
// force the brightness module to be short
func (b *Block) ForceShort() bool {
    return false
}
