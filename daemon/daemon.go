package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/muni-corn/muse-status/format"
)

type Daemon struct {
	addr         string
	leftBlocks   []format.DataBlock
	centerBlocks []format.DataBlock
	rightBlocks  []format.DataBlock
	connections  []net.Conn
}

func New(addr string, leftBlocks, centerBlocks, rightBlocks []format.DataBlock) *Daemon {
	d := &Daemon{
		addr:         addr,
		leftBlocks:   leftBlocks,
		centerBlocks: reverse(centerBlocks),
		rightBlocks:  rightBlocks,
	}

	for _, b := range d.allBlocks() {
		go b.Update()
	}

	return d
}

func (d *Daemon) Start() error {
	// println("starting daemon")

	s, err := net.Listen("tcp", d.addr)
	if err != nil {
		return err
	}

	var currentStatus string

	// accept connections and handle them
	go func() {
		for {
			conn, err := s.Accept()
			if err != nil {
				// println("error on accept:", err.Error())
				continue
			}

			d.handleConnection(conn, currentStatus)
		}
	}()

	outputChan := d.startLemonbarStatus()

	// listen for outputs, and feed them to any connected clients
	go func() {
		for o := range outputChan {
			currentStatus = o
			d.echo(o)
		}
	}()

	go d.listenForXorgChanges()

	return nil
}

func (d *Daemon) HandleCommand(cmd string) error {
	split := strings.SplitN(cmd, " ", 2)

	switch split[0] {
	case "notify":
		if len(split) < 2 {
			break
		}
		notifyName := split[1][:len(split[1])-1] // trims new line
		d.notify(notifyName)
	default:
		return fmt.Errorf("unhandled command: %s", split[0])
	}

	return nil
}

func (d *Daemon) handleConnection(conn net.Conn, init string) {
	d.connections = append(d.connections, conn)

	if format.GetFormatMode() == format.I3JSONMode {
		conn.Write([]byte(`{"version":1}` + "\n["))
	}

	conn.Write([]byte(init + "\n"))

	r := bufio.NewReader(conn)

	go func() {
		str, err := r.ReadString('\n')
		if err != nil {
			return
		}

		// try to handle a command
		err = d.HandleCommand(str)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		}
	}()
}

func (d *Daemon) allBlocks() []format.DataBlock {
	return append(d.leftBlocks, append(d.centerBlocks, d.rightBlocks...)...)
}

func (d *Daemon) startLemonbarStatus() <-chan string {
	// println("starting output")
	outputChannel := make(chan string)

	agg := make(chan bool)
	// for _, b := range centerModules {
	for _, b := range d.allBlocks() {
		c := b.StartBroadcast()
		go func() {
			for v := range c {
				agg <- v
			}
		}()
	}

	go func() {
		// trims first comma on i3bar protocol
		isFirst := true
		for v := range agg {
			if !v {
				continue
			}
			// println("updating status")

			statusStr := d.makeStatusString(&isFirst)

			if format.GetFormatMode() == format.I3JSONMode && isFirst {
				statusStr = statusStr[1:]
				isFirst = true
			}

			outputChannel <- statusStr
		}
	}()

	return outputChannel
}

func (d *Daemon) makeStatusString(isFirst *bool) string {
	switch format.GetFormatMode() {
	case format.I3JSONMode:
		// we will probably want to re-include left modules once we get config files working
		b := format.Chain(append(d.rightBlocks, d.centerBlocks...)...)
		if isFirst != nil && *isFirst {
			*isFirst = false
			return b[1:] // trims comma
		} else {
			return b
		}
	case format.LemonbarMode:
		l := format.Chain(d.leftBlocks...)
		c := format.Chain(d.centerBlocks...)
		r := format.Chain(d.rightBlocks...)
		return fmt.Sprintf("%%{l}%s%%{c}%s%%{r}%s", l, c, r)
	}
	return ""
}

func (d *Daemon) echo(str string) error {
	for _, conn := range d.connections {
		_, err := conn.Write([]byte(str + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemon) echoNewStatus() error {
	status := d.makeStatusString(nil)
	return d.echo(status)
}

func (d *Daemon) notify(what string) {
	for _, b := range d.allBlocks() {
		if b.Name() == what {
			b.Update()
			d.echoNewStatus()
		}
	}
}

func (d *Daemon) listenForXorgChanges() {
	cmd := exec.Command("bspc", "subscribe", "report")
	r, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		// println(err)
		return
	}

	defer r.Close()
	bufr := bufio.NewReader(r)
	for {
		_, _, err := bufr.ReadLine()
		if err != nil {
			// println(err)
		}

		d.notify("bspwm")
		d.notify("window")
	}
}

func reverse(slice []format.DataBlock) []format.DataBlock {
	for left, right := 0, len(slice)-1; left < right; left, right = left+1, right-1 {
		slice[left], slice[right] = slice[right], slice[left]
	}

	return slice
}
