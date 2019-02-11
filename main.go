package main

import (
	"fmt"
	"muse-status/brightness"
	"muse-status/date"
	"muse-status/format"
	// "muse-status/i3"
	"muse-status/mpd"
	"muse-status/network"
	"muse-status/sbattery"
	"muse-status/volume"
	"muse-status/weather"
	// "muse-status/window"
	"os"
)

func main() {
	for k, v := range os.Args {
		if k+1 >= len(os.Args) {
			break
		}
		next := os.Args[k+1]

		switch v {
		case "-S":
			format.SetSecondaryColor(next)
		case "-P":
			format.SetPrimaryColor(next)
		case "-M":
			format.SetFormatMode(next)
		case "-F":
			format.SetTextFont(next)
		case "-I":
			format.SetIconFont(next)
		}
	}

	// channels
	batteryChannel := sbattery.StartSmartBatteryBroadcast()
	brightnessChannel := brightness.StartBrightnessBroadcast()
	dateChannel := date.StartDateBroadcast()
	networkChannel := network.StartNetworkBroadcast()
	volumeChannel := volume.StartVolumeBroadcast()
	weatherChannel := weather.StartWeatherBroadcast()
	mpdChannel := mpd.StartMPDBroadcast()
	// i3Channel := i3.StartI3Broadcast()
	// windowChannel := window.StartWindowBroadcast()

	var battery, brightness, date, mpd, network, volume, weather format.DataBlock

	if format.GetFormatMode() == format.ModeI3Bar {
		fmt.Println(`{ "version": 1 }`)
		fmt.Println("[[]")
	}

	for {
		select {
		case battery = <-batteryChannel:
		case date = <-dateChannel:
		case network = <-networkChannel:
		case volume = <-volumeChannel:
		case brightness = <-brightnessChannel:
		case weather = <-weatherChannel:
		case mpd = <-mpdChannel:
			// case window = <-windowChannel:
			// case i3 = <-i3Channel:
		}

		var status string
		switch format.GetFormatMode() {
		case format.ModeI3Bar:
			status = ",[" + format.Chain(brightness, volume, network, battery, mpd, weather, date) + "]"
		}

		if format.GetFormatMode() == format.ModeLemonbar {
			// add left and right padding
			status = format.Separator() + status
		}

		fmt.Println(status)
	}
}
