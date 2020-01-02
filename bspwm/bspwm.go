package bspwm

import (
	"os/exec"
	"regexp"
	"strings"
	"fmt"

	"github.com/muni-corn/muse-status/format"
)

type workspaceState int

const (
	OccupiedState workspaceState = iota
	ActiveState 
	UrgentState
	DormantState
)

type workspace struct {
	name string
	state workspaceState
}

func getWMStatus() string {
	cmdOutput, err := exec.Command("bspc", "wm", "--get-status").Output()
	if err != nil {
		return ""
	}

	return string(cmdOutput)
}

func parseWorkspaces(wmStatus string) (workspaces []workspace, hasUrgent bool) {
	monitors := regexp.MustCompile("[mM]").Split(wmStatus, -1)[1:] // removes the first element (status prefix), "W"
	for _, m := range monitors {
		// get workspaces
		mWorkspaces := strings.Split(m, ":")
		// monitorName := mWorkspaces[0]
		for i := 1; i < len(mWorkspaces); i++ {
			ws := mWorkspaces[i]
			if len(ws) == 0 {
				continue
			}

			var status workspaceState
			switch ws[0] {
			case 'O', 'F', 'U':
				status = ActiveState
			case 'o':
				status = OccupiedState
			case 'u':
				status = UrgentState
				hasUrgent = true
			default:
				status = DormantState
			}
			newWorkspace := workspace {
				state: status,
				name: ws[1:], // removes the state flag
			}
			workspaces = append(workspaces, newWorkspace)
		}
	}

	return
}

func lemonFormatWorkspaces(ws []workspace) string {
	formatted := ""
	for _, w := range ws {
		var color format.Color
		switch w.state {
		case ActiveState:
			color = format.PrimaryColor()
		case UrgentState:
			color = format.GetWarningColorer().PrimaryColor()
		case DormantState:
			continue
		case OccupiedState:
			fallthrough
		default:
			color = format.SecondaryColor()
		}
		colorStr := color.AlphaHex + color.RGBHex
		formatted += fmt.Sprintf("%%{F#%s}%s%%{F-}    ", colorStr, w.name)
        formatted = fmt.Sprintf("%%{A:bspc desktop -f '^%s':}%s%%{A}", w.name, formatted)
	}

	return strings.TrimSpace(formatted)
}
