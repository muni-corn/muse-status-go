package mpd

import (
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

func getIcon(state playerState) rune {
    switch state {
    case playing:
        return playingIcon
    case paused:
        return pausedIcon
    }

    return playingIcon
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
