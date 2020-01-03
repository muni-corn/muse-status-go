package bspwm

import (
	"github.com/muni-corn/muse-status/format"
)

type Block struct {
	lastStatusOutput    string
	currentStatusOutput string
	hasUrgent           bool

	lemonOutput string

	currentWorkspaces []workspace

	rapidfire bool
}

func NewBSPWMBlock() *Block {
	return new(Block)
}

func (b *Block) StartBroadcast() <-chan bool { // returns a channel that sends signals to update the status bar
	c := make(chan bool)
	return c
}

func (b *Block) Update() {
	b.currentStatusOutput = getWMStatus()
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
