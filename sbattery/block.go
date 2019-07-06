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
	nextStatus ChargeStatus
	nextChargeNow int
}

// NewBlock returns a new sbattery block
func NewBlock(battery string) (*Block, error) {
	b := &Block{battery: battery}

	b.chargeFull, err = utils.GetIntFromFile(b.getBaseDir()+"charge_full")
	if err != nil {
		return nil, err
	}
}

// NextUpdateCheckTime returns the next time at which there should be a check
// for an update
func (b *Block) NextUpdateCheckTime() time.Time {
	return b.nextUpdateTime
}

func (b *Block) NeedsUpdate() bool {
	b.nextChargeNow, err := utils.GetIntFromFile(b.getBaseDir()+"charge_now")
	if err != nil {
		return false
	}
	b.nextStatus, err = b.getBatteryStatus()
	if err != nil {
		return false
	}
	return b.nextChargeNow != b.chargeNow || b.status != b.nextStatus
}

func (b *Block) Update() {
	b.chargeNow = b.nextChargeNow
	b.status = b.nextStatus
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
}

// Colorer returns a colorer depnding on the percentage left on this battery
func (b *Block) Colorer() Colorer {
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

// ForceShort never forces shorting
func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) getBatteryPercentage() int {
	return b.chargeNow * 100 / b.chargeFull
}

func (b *Block) updateBatteryStatus() error {
	str, err := utils.GetStringFromFile(b.getBaseDir()+"status")
	if err != nil {
		return err
	}

	b.status = ChargeStatus(str)
	return nil
}

func (b *Block) getBaseDir() string {
	return SysPowerSupplyBaseDir + "/" + b.battery + "/"
}
