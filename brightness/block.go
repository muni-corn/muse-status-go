package brightness

import (
	"fmt"

	"github.com/muni-corn/muse-status/format"
	"github.com/muni-corn/muse-status/utils"
	"time"
)

const (
	baseDir = "/sys/class/backlight/"
)

var (
	brightnessIcons = [6]rune{'', '', '', '', '', ''}
	// brightnessIcons = [6]rune{'\uf5da', '\uf5db', '\uf5dc', '\uf5dd', '\uf5de', '\uf5df'} // nerd font icons
)

// Block is a block that contains device
// brightness information
type Block struct {
	card              string
	currentBrightness int
	lastBrightness    int
	maxBrightness	  int

	text  string
	icon  rune
	fader *format.FadingColorer
}

// NewBrightnessBlock returns a new brightness.Block
func NewBrightnessBlock(card string) (*Block, error) {
	b := &Block{
		card: card,
	}

	var err error
	b.maxBrightness, err = b.getMaxBrightness()
	if err != nil {
		return nil, err
	}

	b.fader = &format.FadingColorer {
		Duration: 3,
		StartColor: format.PrimaryColor(),
		EndColor: format.SecondaryColor(),
	}

	return b, nil
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
	go b.broadcast(c)
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	b.Update()
	for {
		if b.fader.IsFading() {
			c <- true
		} 
		if b.needsUpdate() {
			println("brightness change")
			b.Update()
			b.fader.Trigger()
			c <- true
		}
		
		time.Sleep(time.Second / 10)
	}
}

// Name returns the name "brightness"
func (b *Block) Name() string {
	return "brightness"
}

// needsUpdate returns true if the Block needs its text
// updated
func (b *Block) needsUpdate() bool {
	var err error
	b.currentBrightness, err = b.getCurrentBrightness()

	if err == nil && b.currentBrightness != b.lastBrightness {
		b.lastBrightness = b.currentBrightness
		return true
	}

	return false
}

// Update updates the text and icon of the block
func (b *Block) Update() {
	b.text = fmt.Sprintf("%d%%", b.currentBrightness*100/b.maxBrightness)
	b.icon = getIcon(b.currentBrightness*100/b.maxBrightness)
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

// Colorer returns a pointer to the block's fader, for color
func (b *Block) Colorer() format.Colorer {
	return b.fader
}

// Hidden always returns false because brightness never has
// a need to be hidden
func (b *Block) Hidden() bool {
	return false
}

// ForceShort returns false because we never really need to
// force the brightness module to be short
func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) getMaxBrightness() (value int, err error) {
	return utils.GetIntFromFile(baseDir + b.card + "/max_brightness")
}

func (b *Block) getCurrentBrightness() (value int, err error) {
	return utils.GetIntFromFile(baseDir + b.card + "/brightness")
}

