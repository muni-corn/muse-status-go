package playerctl

import (
	"github.com/muni-corn/muse-status/format"

	"time"
)

const (
	playingIcon = '\uf387'
	pausedIcon  = '\uf3e4'
)

type Block struct {
	lastTitle, currentTitle   string
	lastArtist, currentArtist string
	lastStatus, currentStatus status

    rapidfire bool
}

func NewPlayerctlBlock(rapidfire bool) *Block {
	if rapidfire {
		println("WARNING! A window block has been created with rapidfire enabled. This can be VERY bad for your system's performance. Try using `muse-status notify volume` instead after volume updates.")
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
		if b.lastTitle != b.currentTitle || b.lastArtist != b.currentArtist || b.lastStatus != b.currentStatus {
			c <- true

			b.lastTitle = b.currentTitle
			b.lastArtist = b.currentArtist
			b.lastStatus = b.currentStatus
		}

		// sleep
		time.Sleep(time.Second/10);
	}
}

func (b *Block) Update() {
	b.currentStatus, _ = getStatus()
	b.currentTitle, _ = getSongTitle()
	b.currentArtist, _ = getArtist()
}

func (b *Block) Name() string {
	return "playerctl"
}

func (b *Block) Hidden() bool {
	return b.currentStatus == stopped || b.currentTitle == ""
}

func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return format.LemonbarOf(b)
}

func (b *Block) Text() (primary, secondary string) {
	return b.currentTitle, b.currentArtist
}

func (b *Block) Icon() rune {
	switch b.currentStatus {
	case playing:
		return playingIcon
	case paused:
		return pausedIcon
	default:
		return ' '
	}
}

func (b *Block) Colorer() format.Colorer {
	return format.GetDefaultColorer()
}
