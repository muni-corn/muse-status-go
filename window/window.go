package window

import (
	"muse-status/format"
	"muse-status/date"
	"strings"
	"os/exec"
	"time"
)

// StartWindowBroadcast returns a string channel that is fed info about the
// current active window. If no window is active, it is fed a greeting or
// information useful to the user.
func StartWindowBroadcast() chan string {
	channel := make(chan string)

	go func() {
		var lastWindow string;
		for {
			currentWindow := window();
			if (lastWindow != currentWindow) {
				channel <- currentWindow;
				lastWindow = currentWindow;
			}

			time.Sleep(time.Second / 10);
		}
	}()

	return channel
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

