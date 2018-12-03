package i3

import (
	"muse-status/format"
	"time"
	"go.i3wm.org/i3"
	"strings"
)

// StartI3Broadcast returns a string channel that is fed
// i3wm workspace and mode information
func StartI3Broadcast() chan string {
	channel := make(chan string)

	go func() {
		eventReceiver := i3.Subscribe(i3.WorkspaceEventType, i3.ModeEventType)
		defer eventReceiver.Close()

		// initialize workspace string
		workspaces, err := i3.GetWorkspaces()
		if err != nil {
			channel <- format.Dim("Couldn't initialize i3 module :(")
		}
		workspacesStr, workspaceUrgency := getWorkspacesString(workspaces)

		var mode string
		var lastMode string
		var lastWorkspacesStr string

		go func() {
			for eventReceiver.Next() {
				event := eventReceiver.Event()

				switch event.(type) {
				case *i3.WorkspaceEvent:
					workspaces, err = i3.GetWorkspaces()
					if err != nil {
						channel <- format.Dim("Couldn't get workspaces")
						continue
					}
					workspacesStr, workspaceUrgency = getWorkspacesString(workspaces)
				case *i3.ModeEvent:
					mode = event.(*i3.ModeEvent).Change
					if mode == "default" {
						mode = ""
					} else {
						if workspacesStr != "" {
							mode = format.Separator() + mode
						}
					}
				}
			}
		}()

		for {
			modeChange := lastMode != mode
			workspacesChange := lastWorkspacesStr != workspacesStr
			if workspaceUrgency || mode != "" || modeChange ||  workspacesChange {
				if workspaceUrgency {
					workspacesStr, workspaceUrgency = getWorkspacesString(workspaces)
					workspacesChange = true
				}
				if modeChange {
					lastMode = mode
				}
				if workspacesChange {
					lastWorkspacesStr = workspacesStr
				}
				channel <- strings.TrimSpace(workspacesStr) + strings.TrimSpace(format.WarningBlink(mode))
			}

			time.Sleep(time.Second / 10)
		}

	}()

	return channel
}

func getWorkspacesString(workspaces []i3.Workspace) (str string, urgency bool) {
	if len(workspaces) <= 1 {
		return
	}

	for k, v := range workspaces {
		if k > 0 {
			str += format.Separator()
		}

		displayStr := v.Name
		if v.Urgent {
			urgency = true
			str += format.WarningBlink(displayStr)
		} else if !v.Focused {
			str += format.Dim(displayStr)
		} else {
			str += displayStr
		}
	}
	return
}
