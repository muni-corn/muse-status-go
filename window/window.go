package window

import (
	"muse-status/format"
	"muse-status/date"
	"strings"
	"os/exec"
	"time"
	"regexp"
)

var (
	lineReturnRegex = regexp.MustCompile(`\r?\n`)
)

// StartWindowBroadcast returns a string channel that is fed info about the
// current active window. If no window is active, it is fed a greeting or
// information useful to the user.
func StartWindowBroadcast() chan *format.ClassicBlock {
	channel := make(chan *format.ClassicBlock)
	block := &format.ClassicBlock{Name: "window", Instance: "window"}

	go func() {
		var lastWindow string;
		for {
			currentWindow := window()
			if (lastWindow != currentWindow) {
				block.PrimaryText = currentWindow
				channel <- block;
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
		return ""
	}

	output := string(cmdOutput)
	if strings.Contains(output, "i3") {
		output = date.GetGreeting()
	} else {
		output = lineReturnRegex.ReplaceAllString(output, "")
	}


	return format.Dim(output)
}

