package mpd

import (
    "muse-status/format"
    "regexp"
    "time"
    "github.com/fhs/gompd/mpd"
)

type playerState int

const (
    stopped playerState = iota
    playing
    paused
)

const (
    playingIcon = '\uf387'
    pausedIcon  = '\uf3e4'
    // nerd font icons
    // playingIcon = '\uf387'
    // pausedIcon  = '\uf38a'
)

var (
    lineSplitRegex = regexp.MustCompile(`\r?\n`)
    statusRegex    = regexp.MustCompile(`\[([A-z]*?)\]`) // matches words within brackets (get first capturing group)
)

// StartMPDBroadcast returns a string channel that is fed information about any
// current media that is playing from an mpd server
func StartMPDBroadcast() chan *format.ClassicBlock {
    channel := make(chan *format.ClassicBlock)
    block := &format.ClassicBlock{Name: "mpd"}

    go func() {
        for {
            // start a client for mpd. if we fail to create one,
            // quit
            mpdClient, err := mpd.Dial("tcp", "localhost:6600")
            if err != nil {
                println("Couldn't start mpd client")
                time.Sleep(time.Second)
                continue
            }

            // create a watcher for mpd; the player subsystem.
            // this will help us to know when changes are made
            // to the current song. if creating the watcher
            // results in an error, we'll display and error and
            // terminate this module
            watcher, err := mpd.NewWatcher("tcp", "localhost:6600", "")
            if err != nil {
                println("Couldn't create mpd watcher")
                time.Sleep(time.Second)
                continue
            }

            if watcher != nil && mpdClient != nil {
                defer watcher.Close()
                defer mpdClient.Close()

                title, artist, state, err := getInfo(mpdClient)
                updateChannel(title, artist, state, block, channel)
                for range watcher.Event {
                    title, artist, state, err = getInfo(mpdClient)
                    if err != nil {
                        // if error, log it
                        println(err.Error())
                        break
                    }

                    updateChannel(title, artist, state, block, channel)
                }
            }

            time.Sleep(time.Second * 2)
        }
    }()

    return channel
}

func updateChannel(title, artist string, state playerState, block *format.ClassicBlock, channel chan *format.ClassicBlock) {
    if state == stopped {
        block.SetHidden(true)
        channel <- block
        return
    } else if block.Hidden() {
        block.SetHidden(false)
    }

    block.PrimaryText = title
    block.SecondaryText = artist

    switch state {
    case playing:
        block.Icon = playingIcon
    case paused:
        block.Icon = pausedIcon
    }

    channel <- block
}

func getInfo(client *mpd.Client) (title string, artist string, state playerState, err error) {
    currentSong, err := client.CurrentSong()
    if err != nil {
        return
    }

    mpdStatus, err := client.Status()
    if err != nil {
        return
    }

    title = currentSong["Title"]
    artist = currentSong["AlbumArtist"]

    rawPlayState := mpdStatus["state"]
    switch rawPlayState {
    case "play":
        state = playing
    case "pause":
        state = paused
    case "stop":
        state = stopped
    }

    return
}
