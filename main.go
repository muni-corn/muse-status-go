package main

import (
	"github.com/muni-corn/muse-status/bspwm"
	"github.com/muni-corn/muse-status/brightness"
	"github.com/muni-corn/muse-status/date"
	"github.com/muni-corn/muse-status/format"
	"github.com/muni-corn/muse-status/network"
	"github.com/muni-corn/muse-status/playerctl"
	"github.com/muni-corn/muse-status/sbattery"
	"github.com/muni-corn/muse-status/volume"
	"github.com/muni-corn/muse-status/weather"
	"github.com/muni-corn/muse-status/window"
	"github.com/muni-corn/muse-status/daemon"

	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const socketAddr = "/tmp/muse-status.sock"

func main() {
	handleArgs()

	bspwmBlock := bspwm.NewBSPWMBlock(false)
	batteryBlock, _ := sbattery.NewSmartBatteryBlock("BAT0", 30, 15)
	brightnessBlock, _ := brightness.NewBrightnessBlock("amdgpu_bl0", false)
	dateBlock := date.NewDateBlock()
	networkBlock, _ := network.NewNetworkBlock("wlo1")
	playerctlBlock := playerctl.NewPlayerctlBlock(false)
	volumeBlock := volume.NewVolumeBlock(false)
	windowBlock := window.NewWindowBlock(false)
	weatherBlock := weather.NewWeatherBlock(nil)

	// TODO parse from configuration file
	d := daemon.New(
		socketAddr, 
		[]format.DataBlock{bspwmBlock, windowBlock},
		[]format.DataBlock{dateBlock, weatherBlock, playerctlBlock},
		[]format.DataBlock{brightnessBlock, volumeBlock, networkBlock, batteryBlock},
	)

	var (
		client net.Conn
		err error
	)

	if client, err = net.Dial("unix", socketAddr); err != nil {
		println("error connecting to socket; starting own daemon")
		err = d.Start()
		if err != nil {
			panic(err)
		}

		client, err = net.Dial("unix", socketAddr)
		if err != nil {
			panic(err)
		}
	}

	handleClient(client)
}

func handleClient(conn net.Conn) error {
	r := bufio.NewReader(conn)

	for {
		str, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Print(str)
	}
}

func handleArgs() {
	// muse be a command if first (second, technically) argument doesn't start
	// with a dash. exit after command
	if len(os.Args) >= 2 && os.Args[1][0] != '-' {
		err := sendCommand(os.Args[1:])
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	for k, v := range os.Args {
		if k+1 >= len(os.Args) {
			break
		}
		next := os.Args[k+1]

		switch v {
		case "-p", "--primary-color":
			format.SetPrimaryColor(next)
		case "-s", "--secondary-color":
			format.SetSecondaryColor(next)
		case "-f", "--font":
			format.SetTextFont(next)
		case "-i", "--icon-font":
			format.SetIconFont(next)
			// case "-r", "--rapid-fire":
			// remove completely? ignore notify actions and "rapid fire" check instead
		}
	}
}

func sendCommand(args []string) error {
	str := strings.Join(args, " ")

	conn, err := net.Dial("unix", socketAddr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(str + "\n"))
	if err != nil {
		panic(err)
	}

	return nil
}

// vim: foldmethod=marker
