package mpd

import (
	"github.com/fhs/gompd/mpd"
	"github.com/muni-corn/muse-status/format"
)

// Block is a format.DataBlock that transmits mpd data
type Block struct {
	addr    string
	watcher *mpd.Watcher
	client  *mpd.Client

	song, artist string
	state        playerState

	needsUpdate bool
	icon        rune
	hidden      bool
}

// NewMPDBlock returns a new Block
func NewMPDBlock(addr string) (b *Block, err error) {
	b = &Block{addr: addr}

	b.watcher, err = mpd.NewWatcher("tcp", addr, "")
	if err != nil {
		return
	}

	b.client, err = mpd.Dial("tcp", addr)

	go func() {
		for range b.watcher.Event {
			b.needsUpdate = true
		}
	}()

	return
}

// NeedsUpdate returns true if an event has happened since
// the last update
func (b *Block) NeedsUpdate() bool {
	return b.needsUpdate
}

// Update updates the mpd block
func (b *Block) Update() {
	title, artist, state, err := getInfo(b.client)
	if err != nil {
		return
	}

	if state == stopped {
		b.hidden = true
		return
	} else if b.Hidden() {
		b.hidden = false
	}

	b.song = title
	b.artist = artist

	switch state {
	case playing:
		b.icon = playingIcon
	case paused:
		b.icon = pausedIcon
	}

	b.needsUpdate = false
}

// Name returns "mpd"
func (b *Block) Name() string {
	return "mpd"
}

// Icon returns a playing or paused icon
func (b *Block) Icon() rune {
	return getIcon(b.state)
}

// Text returns the song as primary, the song artist as secondary
func (b *Block) Text() (primary, secondary string) {
	return b.song, b.artist
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
