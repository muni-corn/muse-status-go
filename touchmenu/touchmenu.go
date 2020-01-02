package touchmenu

import (
    "github.com/muni-corn/muse-status/format"
)

type Block struct{}

func (b *Block) StartBroadcast() <-chan bool {
    return make(chan bool)
}

func (b *Block) Update() { }

func (b *Block) Name() string {
    return "touchmenu"
}

func (b *Block) Hidden() bool {
    return false
}

func (b *Block) ForceShort() bool {
    return false
}

func (b *Block) Output(mode format.Mode) string {
    return "%{A:onboard:}\uf70b%{A}"
}
