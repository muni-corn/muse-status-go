package main

import (
	// "encoding/json"
	"fmt"
	"muse-status/brightness"
	"muse-status/date"
	"muse-status/format"
	"muse-status/network"
	"muse-status/sbattery"
	"muse-status/volume"
	"muse-status/weather"
	"muse-status/window"
	"os/exec"
	"regexp"
	// "go.i3wm.org/i3"
)

func main() {
	// channels
	batteryChannel := sbattery.StartSmartBatteryBroadcast()
	dateChannel := date.StartDateBroadcast()
	networkChannel := network.StartNetworkBroadcast()
	volumeChannel := volume.StartVolumeBroadcast()
	brightnessChannel := brightness.StartBrightnessBroadcast()
	weatherChannel := weather.StartWeatherBroadcast()
	windowChannel := window.StartWindowBroadcast()

	var battery, date, network, volume, brightness, weather, window string

	lineReturnRegex := regexp.MustCompile(`\r?\n`)
	for {
		select {
		case battery = <-batteryChannel:
		case date = <-dateChannel:
		case network = <-networkChannel:
		case volume = <-volumeChannel:
		case brightness = <-brightnessChannel:
		case weather = <-weatherChannel:
		case window = <-windowChannel:
		}

		status := window + format.Center(date+format.Separator()+weather) + " " + format.Right(brightness+format.Separator()+volume+format.Separator()+network+format.Separator()+battery)

		// remove line returns
		status = lineReturnRegex.ReplaceAllString(status, "")

		// add left and right padding
		status = format.Separator() + status + format.Separator()

		fmt.Println(status)
	}
}

func mpd() string {
	output, err := exec.Command("mpc").Output()
	if err != nil {
		return "Error executing mpc. Is it installed?"
	}
	return string(output)

}

type i3Workspace struct {
	num     int
	name    string
	visible bool
	focused bool
	urgent  bool
}
