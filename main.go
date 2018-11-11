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
	"regexp"
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

	lineReturnRegex := regexp.MustCompile(`\r?\n`)
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

		status := window + format.Center(date+format.Separator()+weather+format.Separator()+mpd) + " " + format.Right(brightness+format.Separator()+volume+format.Separator()+network+format.Separator()+battery)

		// remove line returns
		status = lineReturnRegex.ReplaceAllString(status, "")

		// add left and right padding
		status = format.Separator() + status + format.Separator()

		fmt.Println(status)
	}
}
