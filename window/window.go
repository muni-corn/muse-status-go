package window

import (
	"github.com/muni-corn/muse-status/date"

	"os/exec"
	"regexp"
	"strings"
)

var (
	lineReturnRegex = regexp.MustCompile(`\r?\n`)
)

// // StartWindowBroadcast returns a string channel that is fed info about the
// // current active window. If no window is active, it is fed a greeting or
// // information useful to the user.
// func StartWindowBroadcast() chan *format.ClassicBlock {
// 	channel := make(chan *format.ClassicBlock)
// 	block := &format.ClassicBlock{Name: "window"}

// 	go func() {
// 		var lastWindow string
// 		for {
// 			currentWindow := window()
// 			if lastWindow != currentWindow {
// 				block.PrimaryText = currentWindow
// 				channel <- block
// 				lastWindow = currentWindow
// 			}

// 			time.Sleep(time.Second / 10)
// 		}
// 	}()

// 	return channel
// }

// func swayWindow() string {
// 	// get sway tree
// 	cmdOutput, err := exec.Command("swaymsg", "-t", "get_tree").Output()
// 	if err != nil {
// 		return date.GetGreeting()
// 	}

// 	output := string(cmdOutput)
// 	if output == "i3" || strings.TrimSpace(output) == "" {
// 		output = date.GetGreeting()
// 	} else {
// 		output = lineReturnRegex.ReplaceAllString(output, "")
// 	}

// 	return format.Dim(output)
// }

func xWindow() string {
	cmdOutput, err := exec.Command("xdotool", "getwindowfocus", "getwindowname").Output()
	if err != nil {
		return date.GetGreeting()
	}

	window := string(cmdOutput)
	if window == "i3" || strings.TrimSpace(window) == "" {
		window = date.GetGreeting()
	} else {
		window = lineReturnRegex.ReplaceAllString(window, "")
	}

	return window
}
