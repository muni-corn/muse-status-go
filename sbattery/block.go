package sbattery

import (
	"github.com/muni-corn/muse-status/format"
	"github.com/muni-corn/muse-status/utils"
	"time"
	"strconv"
)

type read struct {
	at time.Time
	status ChargeStatus
	charge int
}

// Block is a data block for sbattery
type Block struct {
	warningLevel int
	alarmLevel int

	battery string
	chargeFull int

	lastAnchorRead read // last time the status of the battery changed, used for calculating time remaining
	currentRead read

	nextUpdateTime time.Time
}

// NewSmartBatteryBlock returns a new sbattery block
func NewSmartBatteryBlock(battery string, warningLevel, alarmLevel int) (*Block, error) {
	b := &Block{battery: battery, warningLevel: warningLevel, alarmLevel: alarmLevel}

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
			oldPercentage := b.currentRead.charge / b.chargeFull
			oldStatus := b.currentRead.status

			b.Update()

			newPercentage := b.currentRead.charge / b.chargeFull
			if (b.currentRead.status != oldStatus || newPercentage != oldPercentage) {
				c <- true
			}
		}

		if b.getBatteryPercentage() <= b.warningLevel && b.currentRead.status == Discharging {
			c <- true
			time.Sleep(time.Second / 15)
		} else {
			time.Sleep(b.nextUpdateTime.Sub(time.Now()))
		}
	}
}

func (b *Block) Update() {
	b.nextUpdateTime = time.Now().Add(time.Second * 5)

	newRead, err := b.getNewRead()
	if err != nil {
		return
	}

	b.currentRead = newRead
	if b.currentRead.status != b.lastAnchorRead.status {
		b.lastAnchorRead = b.currentRead
	}
}

// Name returns "battery"
func (b *Block) Name() string {
	return "battery"
}

// Icon returns a battery icon
func (b *Block) Icon() rune {
	return getBatteryIcon(b.currentRead.status, b.getBatteryPercentage())
}

const timeFormat = "3:04 pm"

// Text returns all the text
func (b *Block) Text() (primary, secondary string) {
	primary = strconv.Itoa(b.getBatteryPercentage()) + "%"

	var prefix string
	switch b.currentRead.status {
	case Charging:
		prefix = "Full at"
	case Discharging:
		prefix = "Until"
	default:
		return
	}
	secondary = prefix + " " + b.getCompletionTime().Format(timeFormat)
	return
}

// Colorer returns a colorer depnding on the percentage left on this battery
func (b *Block) Colorer() format.Colorer {
	if b.currentRead.status == Charging {
		return format.GetDefaultColorer()
	}

	perc := b.getBatteryPercentage()
	switch {
	case perc <= b.alarmLevel:
		return format.GetAlarmColorer()
	case perc <= b.warningLevel:
		return format.GetWarningColorer()
	default:
		return format.GetDefaultColorer()
	}
}

// Hidden when status is Full
func (b *Block) Hidden() bool {
	return b.currentRead.status == Full
}

// ForceShort never happens; return false
func (b *Block) ForceShort() bool {
	return false;
}

func (b *Block) Output(mode format.Mode) string {
	return format.LemonbarOf(b)
}

func (b *Block) getNewRead() (read, error) {
	r := read{}

	charge, err := b.getBatteryCharge()
	if err != nil {
		return r, err
	}

	status, err := b.getBatteryStatus()
	if err != nil {
		return r, err
	}

	r.charge = charge
	r.status = status
	r.at = time.Now()

	return r, nil
}

func (b *Block) getBatteryCharge() (int, error) {
	return utils.GetIntFromFile(b.getBaseDir()+"charge_now")
}

func (b *Block) getBatteryPercentage() int {
	return b.currentRead.charge * 100 / b.chargeFull
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

func (b *Block) getCompletionTime() time.Time {
	var end int
	switch b.currentRead.status {
	case Discharging:
		end = 0
	case Charging:
		end = b.chargeFull
	}

	// duration (nanosecond?) per charge unit
	rate := float64(b.currentRead.at.Sub(b.lastAnchorRead.at)) / float64(b.currentRead.charge - b.lastAnchorRead.charge)
	timeLeft := float64(end-b.currentRead.charge) * rate
	return time.Now().Add(time.Duration(timeLeft))
}
