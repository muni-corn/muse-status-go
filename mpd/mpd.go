package mpd

import (
	"errors"
	"muse-status/format"
	"os/exec"
	"regexp"
	"time"
)

type playerState int

const (
	stopped playerState = iota
	playing
	paused
)

const (
	playingIcon = ''
	pausedIcon  = ''
)

var (
	lineSplitRegex = regexp.MustCompile(`\r?\n`)
	statusRegex    = regexp.MustCompile(`\[([A-z]*?)\]`) // matches words within brackets (get first capturing group)
	mpcCmd         = []string{"mpc", "-f", `%title%\n%albumartist%`}
)

// StartMPDBroadcast returns a string channel that is fed information about any
// current media that is playing from an mpd server
func StartMPDBroadcast() chan string {
	channel := make(chan string)

	go func() {
		var lastStatus string

		for {
			title, artist, state, err := getInfo()

			var status string
			if err != nil {
				status = ""
			} else {
				song := title + "  " + format.Dim(artist)

				switch state {
				case playing:
					status = string(playingIcon) + "  " + song
				case paused:
					status = string(pausedIcon) + "  " + song
				case stopped:
					status = ""
				}
			}

			if lastStatus != status {
				channel <- status
				lastStatus = status
			}
			
			time.Sleep(time.Second / 5)
		}
	}()

	return channel
}

func getInfo() (title string, artist string, state playerState, err error) {
	output, err := exec.Command("mpc", "-f", `%title%\n%albumartist%`).Output()
	if err != nil {
		return
	}

	split := lineSplitRegex.Split(string(output), -1)

	if len(split) < 3 {
		err = errors.New("Nothing is in mpd's playlist")
		return
	}

	title = split[0]
	artist = split[1]

	playStateSlice := statusRegex.FindStringSubmatch(split[2])

	if len(playStateSlice) < 2 {
		err = errors.New("There was an issue with MPD")
	}

	rawPlayState := playStateSlice[1] // should be 'on' or 'off'. if it's not, then wtf
	println(rawPlayState)
	switch rawPlayState {
	case "playing":
		state = playing
	case "paused":
		state = paused
	case "stopped":
		state = stopped
	}

	return
}
