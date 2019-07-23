package sbattery

import (
	"github.com/muni-corn/muse-status/format"
	"github.com/muni-corn/muse-status/utils"
	"time"
	"strconv"
)

// Block is a data block for sbattery
type Block struct {
	battery string

	status ChargeStatus
	chargeNow, chargeFull int

	nextUpdateTime time.Time
}

// NewSmartBatteryBlock returns a new sbattery block
func NewSmartBatteryBlock(battery string) (*Block, error) {
	b := &Block{battery: battery}

	var err error
	b.chargeFull, err = utils.GetIntFromFile(b.getBaseDir()+"charge_full")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// StartBroadcast starts broadcasting from this block. It returns a channel
// that sends output when an update should happen
func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
	go b.broadcast(c)
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	for {
		if time.Now().After(b.nextUpdateTime) {
			// store old values
			// use percentage for less aggressive updating
			oldPercentage := b.chargeNow / b.chargeFull
			oldStatus := b.status
			b.Update()

			newPercentage := b.chargeNow / b.chargeFull
			if (b.status != oldStatus || newPercentage != oldPercentage) {
				c <- true
			}
		}

		if b.getBatteryPercentage() <= 30 && b.status == Discharging {
			c <- true
			time.Sleep(time.Second / 10)
		} else {
			time.Sleep(b.nextUpdateTime.Sub(time.Now()))
		}
	}
}

func (b *Block) Update() {
	b.chargeNow, _ = utils.GetIntFromFile(b.getBaseDir()+"charge_now")
	b.status, _ = b.getBatteryStatus()

	b.nextUpdateTime = time.Now().Add(time.Second * 5)
}

// Name returns "battery"
func (b *Block) Name() string {
	return "battery"
}

// Icon returns a battery icon
func (b *Block) Icon() rune {
	return getBatteryIcon(b.status, b.getBatteryPercentage())
}

// Text returns all the text
func (b *Block) Text() (primary, secondary string) {
	primary = strconv.Itoa(b.getBatteryPercentage()) + "%"
	secondary = string(b.status)
	return
}

// Colorer returns a colorer depnding on the percentage left on this battery
func (b *Block) Colorer() format.Colorer {
	if b.status == Charging {
		return format.GetDefaultColorer()
	}

	perc := b.getBatteryPercentage()
	switch {
	case perc <= 15:
		return format.GetAlarmColorer()
	case perc <= 30:
		return format.GetWarningColorer()
	default:
		return format.GetDefaultColorer()
	}
}

// Hidden when status is Full
func (b *Block) Hidden() bool {
	return b.status == Full
}

// ForceShort never happens; return false
func (b *Block) ForceShort() bool {
	return false;
}

func (b *Block) Output(mode format.Mode) string {
	return format.LemonbarOf(b)
}

func (b *Block) getBatteryPercentage() int {
	return b.chargeNow * 100 / b.chargeFull
}

func (b *Block) getBatteryStatus() (ChargeStatus, error) {
	str, err := utils.GetStringFromFile(b.getBaseDir()+"status")
	if err != nil {
		return Unknown, err
	}

	return ChargeStatus(str), nil
}

func (b *Block) getBaseDir() string {
	return SysPowerSupplyBaseDir + "/" + b.battery + "/"
}
