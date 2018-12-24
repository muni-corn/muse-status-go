package network

import (
	"muse-status/format"
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
	signalCmd        = "nmcli -t -f in-use,signal dev wifi | grep '*'"
	statusCmd        = "nmcli -t -f type,state,connection dev"
	connectionIcons  = [5]rune{'冷', '爛', '嵐', '襤', '蠟'}
	disconnectedIcon = '浪'
	disabledIcon     = '來'
)

// StartNetworkBroadcast returns a string channel that is
// fed network information
func StartNetworkBroadcast() chan string {
	channel := make(chan string)

	// TODO should probably write this inline to prevent an accidental
	// synchronous call to broadcast()
	go broadcast(channel)

	return channel
}

// an async function that broadcasts network information to the specified
// channel
func broadcast(channel chan string) {
	for {
		channel <- getWifi()
		time.Sleep(time.Second * updateIntervalSeconds)
	}
}

func getWifi() string {
	// TODO error checking

	output, err := exec.Command("bash", "-c", statusCmd).Output()
	if err != nil {
		println("Error getting connection status");
		return "";
	}

	strOutput := string(output)

	regex := regexp.MustCompile(`\r?\n`)
	interfaces := regex.Split(strOutput, -1)
	var stringToUse string
	for _, v := range interfaces {
		if strings.Contains(v, "wifi") {
			stringToUse = v
			break
		}
	}

	if stringToUse == "" {
		return ""
	}

	if enabled := strings.Contains(stringToUse, "connect"); enabled {
		ssid := strings.Split(stringToUse, ":")[2]
		if strings.Contains(stringToUse, "connecting") {
			// see if we're connecting
			return format.Dim(string(disconnectedIcon) + "  Connecting to " + ssid)
		} else if !strings.Contains(stringToUse, "disconnected") {
			// make sure we're not disconnected from wifi
			signalOuput, err := exec.Command("bash", "-c", signalCmd).Output()
			if err != nil {
				return format.Dim(string(disconnectedIcon))
			}

			// get ssid and signal strength
			ssid := strings.Split(stringToUse, ":")[2]
			signalStr := regex.ReplaceAllString(strings.Split(string(signalOuput), ":")[1], "")
			signal, err := strconv.Atoi(signalStr)

			if err != nil {
				println("Error parsing signal")
				return string(disconnectedIcon) + "  " + ssid
			}

			// get the icon
			iconIndex := len(connectionIcons) * signal / 100

			// constrains index
			if iconIndex >= len(connectionIcons) {
				iconIndex = len(connectionIcons) - 1
			}

			icon := connectionIcons[iconIndex]
			return string(icon) + "  " + ssid
		}
		// if none of the above, we're disconnected
		return format.Dim(string(disconnectedIcon))
	}
	// disabled icon
	return format.Dim(string(disabledIcon))
}

func getEthernet() string {
	return ""
}
