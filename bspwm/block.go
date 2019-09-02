package bspwm

import (
	"github.com/muni-corn/muse-status/format"
	"time"
)

type Block struct {
	lastStatusOutput string
	currentStatusOutput string
	hasUrgent bool

	lemonOutput string

	currentWorkspaces []workspace

    rapidfire bool
}

func NewBSPWMBlock(rapidfire bool) *Block {
	if rapidfire {
		println("WARNING! A bspwm block has been created with rapidfire enabled. This can be VERY bad for your system's performance. Try using `muse-status notify volume` instead after volume updates.")
	}

	return &Block{
        rapidfire: rapidfire,
    }
}

func (b *Block) StartBroadcast() <-chan bool { // returns a channel that sends signals to update the status bar
	c := make(chan bool)
    if b.rapidfire {
        go b.broadcast(c)
    }
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	for {
		if b.hasUrgent {
			b.lemonOutput = lemonFormatWorkspaces(b.currentWorkspaces)
			c <- true
		}
		b.currentStatusOutput = getWMStatus()
		if b.currentStatusOutput != b.lastStatusOutput {
			b.lastStatusOutput = b.currentStatusOutput
			b.Update()
			c <- true
		}

		time.Sleep(time.Second / 10)
	}
}

func (b *Block) Update() {
	b.currentWorkspaces, b.hasUrgent = parseWorkspaces(b.currentStatusOutput)
	b.lemonOutput = lemonFormatWorkspaces(b.currentWorkspaces)
}

func (b *Block) Name() string {
	return "bspwm"
}

func (b *Block) Hidden() bool {
	return false
}

func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return b.lemonOutput
}
