package daemon

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"os"
	"os/signal"
	"syscall"

	"github.com/muni-corn/muse-status/format"
)

type Daemon struct {
	addr string
	leftBlocks []format.DataBlock
	centerBlocks []format.DataBlock
	rightBlocks []format.DataBlock
	connections []net.Conn
}

func New(addr string, leftBlocks, centerBlocks, rightBlocks []format.DataBlock) *Daemon {
	d := &Daemon{
		addr: addr, 
		leftBlocks: leftBlocks, 
		centerBlocks: centerBlocks, 
		rightBlocks: rightBlocks,
	}

	for _, b := range d.allBlocks() {
		b.Update()
	}

	return d
}

func (d *Daemon) Start() error {
	println("starting daemon")

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
				println("error on accept:", err.Error())
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

	return nil
}

func (d *Daemon) HandleCommand(cmd string) error {
	println("handling command:", cmd)

	split := strings.SplitN(cmd, " ", 2)
	
	switch split[0] {
	case "notify":
		if len(split) < 2 {
			break;
		}
		notifyName := split[1][:len(split[1])-1] // trims new line
		println("notifying \"" + notifyName + "\"")
		d.notify(notifyName)
	default:
		return fmt.Errorf("unhandled command: %s", split[0])
	}

	return nil
}

func (d *Daemon) handleConnection(conn net.Conn, init string) {
	d.connections = append(d.connections, conn)

	conn.Write([]byte(init + "\n"))

	r := bufio.NewReader(conn)

	go func() {
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				return
			}

			err = d.HandleCommand(str)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			}
		}
	}()
}

func (d *Daemon) allBlocks() []format.DataBlock {
	return append(d.leftBlocks, append(d.centerBlocks, d.rightBlocks...)...)
}

func (d *Daemon) startLemonbarStatus() <-chan string {
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
		for v := range agg {
			if !v {
				continue
			}

			outputChannel <- d.makeStatusString()
		}
	}()

	return outputChannel
}

func (d *Daemon) makeStatusString() string {
	l := format.Chain(d.leftBlocks...)
	c := format.Chain(d.centerBlocks...)
	r := format.Chain(d.rightBlocks...)
	return fmt.Sprintf("%%{l}%s%%{c}%s%%{r}%s", l, c, r)
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
	status := d.makeStatusString()
	return d.echo(status)
}

func (d *Daemon) notify(what string) {
	for _, b := range d.allBlocks() {
		print("notify ", b.Name(), "? ")
		if b.Name() == what {
			println("yes")
			b.Update()
			d.echoNewStatus()
		} else {
			println("no")
		}
	}
}
