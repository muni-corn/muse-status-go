package network

import (
	"github.com/mdlayher/wifi"
	"github.com/muni-corn/muse-status/format"

	"time"
	"errors"
)

// Block is a block that transmits time and date data
type Block struct {
	iface *wifi.Interface
	client *wifi.Client

	currentSSID string
	currentStrengthPct int
	currentStatus	  networkStatus

	lastSSID string
	lastStrengthPct int
	lastStatus	  networkStatus
}

// NewDateBlock returns a new date.Block
func NewNetworkBlock(interfaceName string) (*Block, error) {
	client, err := wifi.New()
	if err != nil {
		return nil, err
	}

	// get all interfaces
	ifs, err := client.Interfaces()
	if err != nil {
		return nil, err
	}
	iface := getInterface(interfaceName, ifs)
	if iface == nil {
		return nil, errors.New("no interface found for " + interfaceName)
	}

	// but only select the one we want
	return &Block{
		client: client,
		iface: iface,
	}, nil
}

// only returns one Interface that matches the name given
func getInterface(interfaceName string, interfaces []*wifi.Interface) *wifi.Interface {
	for _, i := range interfaces {
		if i.Name == interfaceName {
			return i
		}
	}

	return nil
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
	go b.broadcast(c)
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	for {
		b.Update()
		if b.shouldNotify() {
			c <- true
		}

		time.Sleep(time.Second * 5)
	}
}

func (b *Block) shouldNotify() bool {
	if b.lastSSID != b.currentSSID || b.lastStatus != b.currentStatus || b.lastStrengthPct != b.currentStrengthPct {
		b.lastSSID = b.currentSSID
		b.lastStatus = b.currentStatus
		b.lastStrengthPct = b.currentStrengthPct
		return true
	}
	return false
}

// Update updates the network information
func (b *Block) Update() {
	// strength
	infos, err := b.client.StationInfo(b.iface)
	if err != nil {
		return
	}

	if len(infos) > 0 {
		// only going to worry about the first
		dbm := infos[0].Signal
		b.currentStrengthPct = int(dBmToPercentage(float32(dbm)))
	} else {
		b.currentStatus = disconnectedStatus
		return
	}

	// ssid
	bss, err := b.client.BSS(b.iface)
	if err != nil {
		b.currentStatus = disconnectedStatus
		return
	}

	if packetLoss(b.iface.Name) {
		b.currentStatus = packetLossStatus
	} else {
		b.currentStatus = connectedStatus
	}
	b.currentSSID = bss.SSID
}

const (
	signalMaxDBm = -20
	noiseFloorDBm = -90
)

// thank u to i3status and NetworkManager :)
func dBmToPercentage(dbm float32) float32 {
	if (dbm < noiseFloorDBm) {
		dbm = noiseFloorDBm
	}
	if (dbm > signalMaxDBm) {
		dbm = signalMaxDBm
	}

	return -0.008*dbm*dbm + 0.2*dbm + 100
}

// Name returns "network"
func (b *Block) Name() string {
	return "network"
}

// Icon returns the network icon
func (b *Block) Icon() rune {
	return getIcon(b.currentStrengthPct, b.currentStatus)
}

// Text returns the ssid as primary, the status as secondary
func (b *Block) Text() (primary, secondary string) {
	return b.currentSSID, string(b.currentStatus)
}

// Colorer returns the default colorer
func (b *Block) Colorer() format.Colorer {
	return format.GetDefaultColorer()
}

// Hidden returns true if disconnected
func (b *Block) Hidden() bool {
	return b.currentStatus == disconnectedStatus
}

// ForceShort returns false; no force-shorting date yet
func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return format.LemonbarOf(b)
}

func getIcon(signalStrengthPct int, status networkStatus) rune {
	// determine which icons we'll use based on
	// packetLoss
	var icons []rune
	if status == packetLossStatus {
		icons = packetLossIcons
	} else {
		icons = connectionIcons
	}

	// get the icon
	iconIndex := len(icons) * signalStrengthPct / 100

	// constrains index
	if iconIndex >= len(icons) {
		iconIndex = len(icons) - 1
	}

	return icons[iconIndex]
}
