package network

import (
	"github.com/muni-corn/muse-status/format"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	updateIntervalSeconds = 5 // interval to update network information, in seconds
)

var (
	signalCmd         = "nmcli -t -f in-use,signal dev wifi | grep '*'"
	statusCmd         = "nmcli -t -f type,state,connection dev"
	wirelessInterface = "wlo1"
	connectionIcons   = []rune{'\uf92e', '\uf91e', '\uf921', '\uf924', '\uf927'}
	packetLossIcons   = []rune{'\uf92a', '\uf91f', '\uf922', '\uf925', '\uf928'}
	disconnectedIcon  = '\uf92e'
	disabledIcon      = '\uf92d'
)

// StartNetworkBroadcast returns a string channel that is
// fed network information
func StartNetworkBroadcast() chan *format.ClassicBlock {
	channel := make(chan *format.ClassicBlock)
	block := &format.ClassicBlock{Name: "network"}

	// TODO should probably write this inline to prevent an accidental
	// synchronous call to broadcast()
	go broadcast(block, channel)

	return channel
}

// an async function that broadcasts network information to the specified
// channel
func broadcast(block *format.ClassicBlock, channel chan *format.ClassicBlock) {
	for {
		getWifi(block)
		channel <- block
		time.Sleep(time.Second * updateIntervalSeconds)
	}
}

var lineReturnRegex = regexp.MustCompile(`\r?\n`)

func getWifi(block *format.ClassicBlock) {
	// TODO error checking

	output, err := exec.Command("bash", "-c", statusCmd).Output()
	if err != nil {
		println("Error getting connection status")
		return
	}

	strOutput := string(output)

	interfaces := lineReturnRegex.Split(strOutput, -1)
	var stringToUse string
	for _, v := range interfaces {
		if strings.Contains(v, "wifi") {
			stringToUse = v
			break
		}
	}

	if stringToUse == "" {
		return
	}

	if enabled := strings.Contains(stringToUse, "connect"); enabled {
		ssid := strings.Split(stringToUse, ":")[2]
		if strings.Contains(stringToUse, "connecting") {
			// see if we're connecting
			block.Set(format.UrgencyLow, disconnectedIcon, "Connecting to "+ssid, "")
			return
		} else if !strings.Contains(stringToUse, "disconnected") { // make sure we're not disconnected from wifi
			ssid := strings.Split(stringToUse, ":")[2]
			signal, err := getSignalStrength(block)

			if err != nil {
				block.Set(format.UrgencyLow, disconnectedIcon, ssid, "")
				return
			}

			// determine which icons we'll use based on
			// packetLoss
			var icons []rune
			if packetLoss() {
				icons = packetLossIcons
			} else {
				icons = connectionIcons
			}

			// get the icon
			iconIndex := len(icons) * signal / 100

			// constrains index
			if iconIndex >= len(icons) {
				iconIndex = len(icons) - 1
			}

			block.Set(format.UrgencyNormal, icons[iconIndex], ssid, "")
			return
		}
		// if none of the above, we're disconnected
		block.Set(format.UrgencyLow, disconnectedIcon, "", "")
		return
	}
	// disabled icon
	block.Set(format.UrgencyLow, disabledIcon, "", "")
}

func getSignalStrength(block *format.ClassicBlock) (signal int, err error) {
	signalOuput, err := exec.Command("bash", "-c", signalCmd).Output()
	if err != nil {
		block.Urgency = format.UrgencyLow
		block.Icon = disconnectedIcon
		block.PrimaryText = ""
		return
	}

	// get ssid and signal strength
	signalStr := lineReturnRegex.ReplaceAllString(strings.Split(string(signalOuput), ":")[1], "")
	signal, err = strconv.Atoi(signalStr)
	return
}

func getEthernet() string {
	return ""
}

func packetLoss() bool {
	cmd := exec.Command("ping", "-c", "2", "-W", "2", "-I", wirelessInterface, "8.8.8.8")
	err := cmd.Run()
	return err != nil
}
