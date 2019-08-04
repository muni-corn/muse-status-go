package main

import (
	"fmt"
	"github.com/muni-corn/muse-status/bspwm"
	"github.com/muni-corn/muse-status/brightness"
	"github.com/muni-corn/muse-status/date"
	"github.com/muni-corn/muse-status/format"
	// "github.com/muni-corn/muse-status/mpd"
	// "github.com/muni-corn/muse-status/network"
	"github.com/muni-corn/muse-status/playerctl"
	"github.com/muni-corn/muse-status/sbattery"
	"github.com/muni-corn/muse-status/volume"
	"github.com/muni-corn/muse-status/weather"
	"github.com/muni-corn/muse-status/window"
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
	bspwmBlock := bspwm.NewBSPWMBlock()
	batteryBlock, _ := sbattery.NewSmartBatteryBlock("BAT0")
	brightnessBlock, _ := brightness.NewBrightnessBlock("amdgpu_bl0")
	dateBlock := date.NewDateBlock()
	playerctlBlock := playerctl.NewPlayerctlBlock()
	volumeBlock := volume.NewVolumeBlock()
	windowBlock := window.NewWindowBlock()
	weatherBlock := weather.NewWeatherBlock(nil)

	dateChan := dateBlock.StartBroadcast()
	playerctlChan := playerctlBlock.StartBroadcast()
	brightnessChan := brightnessBlock.StartBroadcast()
	batteryChan := batteryBlock.StartBroadcast()
	volumeChan := volumeBlock.StartBroadcast()
	windowChan := windowBlock.StartBroadcast()
	bspwmChan := bspwmBlock.StartBroadcast()
	weatherChan := weatherBlock.StartBroadcast()

	leftModules := []format.DataBlock{bspwmBlock, windowBlock};
	middleModules := []format.DataBlock{dateBlock, weatherBlock, playerctlBlock};
	rightModules := []format.DataBlock{brightnessBlock, volumeBlock, batteryBlock};

	for {
		// hold until an update
		select {
		case <-dateChan:
		case <-brightnessChan:
		case <-batteryChan:
		case <-volumeChan:
		case <-windowChan:
		case <-bspwmChan:
		case <-playerctlChan:
		case <-weatherChan:
		}

		l := format.Chain(leftModules...)
		c := format.Chain(middleModules...)
		r := format.Chain(rightModules...)
		fmt.Printf("%%{l}%s%%{c}%s%%{r}%s\n", l, c, r)
	}
}

// vim: foldmethod=marker
