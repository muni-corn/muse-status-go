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
	"os/exec"
	"regexp"
	"strings"
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

	var battery string
	var date string
	var network string
	var volume string
	var brightness string
	var weather string
	var mpd string

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
		}

		status := window() + format.Center(date+format.Separator()+weather+format.Separator()+mpd) + " " + format.Right(brightness+format.Separator()+volume+format.Separator()+network+format.Separator()+battery)

		// remove line returns
		status = lineReturnRegex.ReplaceAllString(status, "")

		// add left and right padding
		status = format.Separator() + status + format.Separator()

		fmt.Println(status)
	}
}

func window() string {
	cmdOutput, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	if err != nil {
		return "Error executing xdotool. Is it installed?"
	}

	output := string(cmdOutput)
	if strings.Contains(output, "i3") {
		output = date.GetGreeting()
	}

	return format.Dim(output)
}
