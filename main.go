package main

import (
	"encoding/json"
	"fmt"
	"muse-status/date"
	f "muse-status/format"
	bat "muse-status/sbattery"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	// channels
	batteryChannel := bat.StartSmartBatteryBroadcast()
	dateChannel := date.StartDateBroadcast()

	var battery string
	var date string

	lineReturnRegex := regexp.MustCompile(`\r?\n`)
	for {
		select {
		case battery = <-batteryChannel:
		case date = <-dateChannel:
		}

		status := window() + f.Center(date) + " " + f.Right(battery)

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

	return f.Dim(output)
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

func i3() string {
	output, err := exec.Command("i3-msg", "-t", "get_workspaces").Output()
	if err != nil {
		return "Error executing i3-msg."
	}
	var workspaces []i3Workspace
	err = json.Unmarshal(output, &workspaces)
	if err != nil {
		return "Couldn't process i3 json: " + err.Error()
	}
	// TODO
	return ""
}
