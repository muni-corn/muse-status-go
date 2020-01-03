package weather

import (
	"github.com/muni-corn/muse-status/format"

	"time"
)

type Block struct {
	currentReport fullWeatherReport

	location *WeatherLocation
}

func NewWeatherBlock(loc *WeatherLocation) *Block {
	if loc == nil {
		loc, _ = getLocation()
	}
	return &Block{location: loc}
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
	go b.broadcast(c)
	return c
}

func (b *Block) broadcast(c chan<- bool) {
	for {
		b.Update()
		c <- true
		time.Sleep(time.Minute * 20)
	}
}

func (b *Block) Update() {
	b.currentReport, _ = getFullWeatherReport(b.location)
}

func (b *Block) Name() string {
	return "weather"
}

func (b *Block) Hidden() bool {
	return false
}

func (b *Block) ForceShort() bool {
	return false
}

func (b *Block) Output(mode format.Mode) string {
	return format.FormatClassicBlock(b)
}

func (b *Block) Text() (primary, secondary string) {
	return getTemperatureString(b.currentReport), getWeatherDescription(b.currentReport)
}

func (b *Block) Icon() rune {
	return getWeatherIcon(b.currentReport)
}

func (b *Block) Colorer() format.Colorer {
	return format.GetDefaultColorer()
}
