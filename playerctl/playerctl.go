package playerctl

import (
	"os/exec"
	"bytes"
)

type status string

const (
	playing 	status = "Playing"
	paused 			   = "Paused"
	stopped 		   = "Stopped"
)

func getSongTitle() (string, error) {
	return getMetadata("title")
}

func getArtist() (string, error) {
	return getMetadata("artist")
}

func getMetadata(name string) (string, error) {
	cmdOutput, err := exec.Command("playerctl", "metadata", name).Output()
	if err != nil {
		return "error", err
	}

	return string(bytes.Trim(cmdOutput, "\n\r ")), nil
}

func getStatus() (status, error) {
	cmdOutput, err := exec.Command("playerctl", "status").Output()
	if err != nil {
		return stopped, err
	}

	statusRaw := string(bytes.Trim(cmdOutput, "\n\r "))
	switch statusRaw {
	case "Playing":
		return playing, nil
	case "Paused":
		return paused, nil
	default:
		return stopped, nil
	}
}
