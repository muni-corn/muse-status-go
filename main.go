package main

import (
	"fmt"
	"muse-status/brightness"
	"muse-status/date"
	"muse-status/format"
	"muse-status/mpd"
	"muse-status/network"
	"muse-status/sbattery"
	"muse-status/volume"
	"muse-status/weather"
	"muse-status/window"
)

func main() {
	// channels
	batteryChannel := sbattery.StartSmartBatteryBroadcast()
	dateChannel := date.StartDateBroadcast()
	networkChannel := network.StartNetworkBroadcast()
	volumeChannel := volume.StartVolumeBroadcast()
	brightnessChannel := brightness.StartBrightnessBroadcast()
	weatherChannel := weather.StartWeatherBroadcast()
	mpdChannel := mpd.StartMPDBroadcast()
	windowChannel := window.StartWindowBroadcast()

	var battery, date, mpd, network, volume, brightness, weather, window string

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
		}

		status := window + format.Center(format.Chain(date, weather, mpd)) + format.Right(format.Chain(brightness, volume, network, battery))

		// add left and right padding
		status = format.Separator() + status

		fmt.Println(status)
	}
}
