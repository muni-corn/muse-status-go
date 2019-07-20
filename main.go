package main

import (
	"fmt"
	"github.com/muni-corn/muse-status/brightness"
	"github.com/muni-corn/muse-status/date"
	"github.com/muni-corn/muse-status/format"
	// "github.com/muni-corn/muse-status/mpd"
	// "github.com/muni-corn/muse-status/network"
	"github.com/muni-corn/muse-status/sbattery"
	// "github.com/muni-corn/muse-status/volume"
	// "github.com/muni-corn/muse-status/weather"
	// "github.com/muni-corn/muse-status/window"
	"os"
	// "time"
)

func main() {
	for k, v := range os.Args {
		if k+1 >= len(os.Args) {
			break
		}
		next := os.Args[k+1]

		switch v {
		case "-lemonbar":
			format.SetFormatMode(format.LemonbarMode);
		case "-S":
			format.SetSecondaryColor(next)
		case "-P":
			format.SetPrimaryColor(next)
		case "-F":
			format.SetTextFont(next)
		case "-I":
			format.SetIconFont(next)
		}
	}

	switch (format.GetFormatMode()) {
	case format.LemonbarMode:
		lemonbarStatus()
	default:
		lemonbarStatus()
	}
}

func lemonbarStatus() {
	batteryBlock, _ := sbattery.NewSmartBatteryBlock("BAT0")
	brightnessBlock, _ := brightness.NewBrightnessBlock("amdgpu_bl0")
	dateBlock := date.NewDateBlock()

	dateChan := dateBlock.StartBroadcast()
	brightnessChan := brightnessBlock.StartBroadcast()
	batteryChan := batteryBlock.StartBroadcast()

	leftModules := []format.DataBlock{};
	middleModules := []format.DataBlock{dateBlock};
	rightModules := []format.DataBlock{brightnessBlock, batteryBlock};

	for {
		// hold until an update
		select {
		case <-dateChan:
		case <-brightnessChan:
		case <-batteryChan:
		}

		l := format.Chain(leftModules...)
		c := format.Chain(middleModules...)
		r := format.Chain(rightModules...)
		fmt.Printf("%%{l}%s%%{c}%s%%{r}%s\n", l, c, r)
	}
}

// func lemonRunModules(blocks ...format.DataBlock) <-chan string { {{{
// 	for _, b := range blocks {
// 		b.Update()
// 	}

// 	ch := make(chan string)

// 	go func() {
// 		for {
// 			var nextUpdate time.Time

// 			// check for when next check should be
// 			for _, b := range blocks {
// 				bNextUpdateTime := b.NextUpdate()
// 				if bNextUpdateTime.Before(nextUpdate) {
// 					nextUpdate = bNextUpdateTime
// 				}
// 			}

// 			// sleep
// 			time.Sleep(nextUpdate.Sub(time.Now()));

// 			// update blocks
// 			var printNeeded bool
// 			for _, b := range blocks {
// 				if time.Now().After(b.NextUpdate()) {
// 					b.Update()
// 					printNeeded = true
// 				}
// 			}

// 			if updateNeeded {
// 				ch <- format.Chain(blocks...)
// 			}
// 		}
// 	}()

// 	return ch
// }

// if format.GetFormatMode() == format.ModeI3Bar {
// 	fmt.Println(`{ "version": 1 }`)
// 	fmt.Println("[[]")
// }

// for {
// 	select {
// 	case battery = <-batteryChannel:
// 	case date = <-dateChannel:
// 	case network = <-networkChannel:
// 	case volume = <-volumeChannel:
// 	case brightness = <-brightnessChannel:
// 	case weather = <-weatherChannel:
// 	case mpd = <-mpdChannel:
// 		// case window = <-windowChannel:
// 	}

// 	var status string
// 	switch format.GetFormatMode() {
// 	case format.ModeI3Bar:
// 		status = ",[" + format.Chain(brightness, volume, network, battery, mpd, weather, date) + "]"
// 	}

// 	if format.GetFormatMode() == format.ModeLemonbar {
// 		// add left and right padding
// 		status = format.Separator() + status
// 	}

// 	fmt.Println(status)
// }
// } }}}

// vim: foldmethod=marker
