package main

import (
	"fmt"
	"muse-status/brightness"
	"muse-status/date"
	"muse-status/format"
	"muse-status/i3"
	"muse-status/mpd"
	"muse-status/network"
	"muse-status/sbattery"
	"muse-status/volume"
	"muse-status/weather"
	"muse-status/window"
	"os"
)

func main() {
	for k, v := range os.Args {
		if v == "-S" {
			format.SetSecondaryColor(os.Args[k+1])
		}
	}

	// channels
	batteryChannel := sbattery.StartSmartBatteryBroadcast()
	dateChannel := date.StartDateBroadcast()
	networkChannel := network.StartNetworkBroadcast()
	volumeChannel := volume.StartVolumeBroadcast()
	brightnessChannel := brightness.StartBrightnessBroadcast()
	weatherChannel := weather.StartWeatherBroadcast()
	mpdChannel := mpd.StartMPDBroadcast()
	i3Channel := i3.StartI3Broadcast()
	windowChannel := window.StartWindowBroadcast()

	var battery, date, i3, mpd, network, volume, brightness, weather, window string

	for {
		select {
		case battery = <-batteryChannel:
		case date = <-dateChannel:
		case network = <-networkChannel:
		case volume = <-volumeChannel:
		case brightness = <-brightnessChannel:
		case weather = <-weatherChannel:
		case mpd = <-mpdChannel:
		case window = <-windowChannel:
		case i3 = <-i3Channel:
		}

		status := format.Chain(i3, window) +
			format.Center(format.Chain(date, weather, mpd)) +
			format.Right(format.Chain(brightness, volume, network, battery))

		// add left and right padding
		status = format.Separator() + status

		fmt.Print(status, "\r")
	}
}
