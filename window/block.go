package window

import (
	"github.com/muni-corn/muse-status/format"

	"time"
)

type Block struct {
	currentWindow string
	lastWindow string
    rapidfire bool
}

func NewWindowBlock(rapidfire bool) *Block {
	if rapidfire {
		println("WARNING! A window block has been created with rapidfire enabled. This can be VERY bad for your system's performance. Try using `muse-status notify window` instead after window updates.")
	}

	return &Block{
		rapidfire: rapidfire,
	}
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
    if b.rapidfire {
        go b.broadcast(c)
    }
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	for {
		b.Update()
		if b.currentWindow != b.lastWindow {
			c <- true
			b.lastWindow = b.currentWindow
		}
		time.Sleep(time.Second / 10)
	}
}

func (b *Block) Update() {
	b.currentWindow = xWindow()
}

func (b *Block) Name() string {
	return "window"
}

func (b *Block) Icon() rune {
	return ' '
}

func (b *Block) Text() (primary, secondary string) {
	return "", b.currentWindow
}

func (b *Block) Colorer() format.Colorer {
	return format.GetDefaultColorer()
}

func (b *Block) Hidden() bool {
	return false
}

func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return format.LemonbarOf(b)
}
