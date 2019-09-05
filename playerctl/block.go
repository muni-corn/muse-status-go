package playerctl

import (
	"github.com/muni-corn/muse-status/format"
    "bufio"
    "os/exec"
)

const (
	playingIcon = '\uf387'
	pausedIcon  = '\uf3e4'
)

type Block struct {
	lastTitle, currentTitle   string
	lastArtist, currentArtist string
	lastStatus, currentStatus status
}

func NewPlayerctlBlock() *Block {
	return &Block{}
}

func (b *Block) StartBroadcast() <-chan bool {
	c := make(chan bool)
    go b.broadcast(c)
	return c
}

func (b *Block) broadcast(c chan<- bool) {
    metaCmd := exec.Command("playerctl", "metadata", "--follow")
    statusCmd := exec.Command("playerctl", "status", "--follow")
    mr, err := metaCmd.StdoutPipe()
    if err != nil {
        return
    }
    sr, err := statusCmd.StdoutPipe()
    if err != nil {
        return
    }

    metaCmdReader := bufio.NewReader(mr)
    statusCmdReader := bufio.NewReader(sr)

    go func() {
        err = metaCmd.Start()
        if err != nil {
            println(err)
            return
        }

        defer mr.Close()
        for {
            in, _, err := metaCmdReader.ReadLine()
            if err != nil {
                println(err)
                continue
            }

            println(string(in))
            b.Update()
            c <- true
        }
    }()

    go func() {
        err = statusCmd.Start()
        if err != nil {
            println(err)
            return
        }

        defer sr.Close()
        for {
            _, _, err := statusCmdReader.ReadLine()
            if err != nil {
                println(err)
                continue
            }

            b.Update()
            c <- true
        }
    }()
}

func (b *Block) Update() {
    b.currentStatus, _ = getStatus()
    b.currentTitle, _ = getSongTitle()
    b.currentArtist, _ = getArtist()
}

func (b *Block) Name() string {
    return "playerctl"
}

func (b *Block) Hidden() bool {
    return b.currentStatus == stopped || b.currentTitle == ""
}

func (b *Block) ForceShort() bool {
    return false
}

func (b *Block) Output(mode format.Mode) string {
    return format.LemonbarOf(b)
}

func (b *Block) Text() (primary, secondary string) {
    return b.currentTitle, b.currentArtist
}

func (b *Block) Icon() rune {
    switch b.currentStatus {
    case playing:
        return playingIcon
    case paused:
        return pausedIcon
    default:
        return ' '
    }
}

func (b *Block) Colorer() format.Colorer {
    return format.GetDefaultColorer()
}
