package main

import (
	"fmt"
	"os/exec"
	"time"
	"regexp"
	f "muse-status/format"
	bat "muse-status/sbattery"
	"encoding/json"
)

func main() {
	regex := regexp.MustCompile(`\r?\n`);
	for {
		status := window() + f.Center(date()) + " " + f.Right(bat.Status())

		// add background and remove line returns
		status = regex.ReplaceAllString(status, "")

		// add left and right padding
		status = "        " + status + "        "

		fmt.Println(status);
		time.Sleep(time.Second / 2)
	}
}

func battery() string {
	output, err := exec.Command("acpi").Output()
	if err != nil {
		return "Error executing acpi. Is it installed?"
	}
	return string(output)
}

func window() string {
	output, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	if err != nil {
		return "Error executing xdotool. Is it installed?"
	}
	return f.Dim(string(output))
}

func mpd() string {
	output, err := exec.Command("mpc").Output()
	if err != nil {
		return "Error executing mpc. Is it installed?"
	}
	return string(output)

}

type i3Workspace struct {
	num int
	name string
	visible bool
	focused bool
	urgent bool
}

func i3() string {
	output, err := exec.Command("i3-msg", "-t", "get_workspaces").Output();
	if err != nil {
		return "Error executing i3-msg."
	}
	var workspaces []i3Workspace
	err = json.Unmarshal(output, &workspaces)
	if (err != nil) {
		return "Couldn't process i3 json: " + err.Error();
	}
	// TODO
	return ""
}
