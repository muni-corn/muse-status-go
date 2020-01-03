package volume

import (
	"fmt"
	"time"

	"github.com/muni-corn/muse-status/format"
)

type Block struct {
	lastVolume    int
	currentVolume int
	rapidfire     bool

	fader *format.FadingColorer
}

func NewVolumeBlock(rapidfire bool) *Block {
	if rapidfire {
		// println("WARNING! A volume block has been created with rapidfire enabled. This can be VERY bad for your system's performance. Try using `muse-status notify volume` instead after volume updates.")
	}

	b := &Block{
		rapidfire: rapidfire,
	}
	b.fader = &format.FadingColorer{
		Duration:   3,
		StartColor: format.PrimaryColor(),
		EndColor:   format.SecondaryColor(),
	}
	return b
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
	if b.rapidfire {
		go b.broadcast(c)
	}
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	// init
	b.Update()
	c <- true

	for b.rapidfire {
		b.Update()

		if b.fader.IsFading() {
			c <- true
		}
		if b.currentVolume != b.lastVolume {
			c <- true
			b.lastVolume = b.currentVolume
			b.fader.Trigger()
		}

		time.Sleep(time.Second / 10)
	}
}

func (b *Block) Update() {
	b.currentVolume, _ = getCurrentVolume()
}

func (b *Block) Name() string {
	return "volume"
}

func (b *Block) Icon() rune {
	return getIcon(b.currentVolume)
}

func (b *Block) Text() (primary, secondary string) {
	if b.currentVolume > 0 {
		return fmt.Sprintf("%d%%", b.currentVolume), ""
	} else {
		return "Muted", ""
	}
}

func (b *Block) Colorer() format.Colorer {
	return b.fader
}

func (b *Block) Hidden() bool {
	return false
}

func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return format.FormatClassicBlock(b)
}
