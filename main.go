package main

import (
	// "encoding/json"
	"fmt"
	"muse-status/brightness"
	"muse-status/date"
	"muse-status/format"
	"muse-status/sbattery"
	"muse-status/network"
	// "muse-status/volume"
	"os/exec"
	"regexp"
	"strings"
	// "go.i3wm.org/i3"
)

func main() {
	// channels
	batteryChannel := sbattery.StartSmartBatteryBroadcast()
	dateChannel := date.StartDateBroadcast()
	networkChannel := network.StartNetworkBroadcast()
	// volumeChannel := volume.StartVolumeBroadcast()
	brightnessChannel := brightness.StartBrightnessBroadcast()

	var battery string
	var date string
	var network string
	// var volume string
	var brightness string

	lineReturnRegex := regexp.MustCompile(`\r?\n`)
	for {
		select {
		case battery = <-batteryChannel:
		case date = <-dateChannel:
		case network = <-networkChannel:
		// case volume = <-volumeChannel:
		case brightness = <-brightnessChannel:
		}

		status := window() + format.Center(date) + " " + format.Right(brightness + format.Separator() + network + format.Separator() + battery)

		// remove line returns
		status = lineReturnRegex.ReplaceAllString(status, "")

		// add left and right padding
		status = "        " + status + "        "

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

