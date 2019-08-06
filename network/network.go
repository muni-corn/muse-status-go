package network

import (
	// "github.com/muni-corn/muse-status/format"
	"os/exec"
	"regexp"
	// "strconv"
	// "strings"
)

type networkStatus string

const (
	disconnectedStatus networkStatus = "No connection"
	packetLossStatus                 = "Can't reach the Internet"
	connectingStatus                 = "Connecting"
	connectedStatus                  = ""
	signInRequired                   = "Sign-in required"
	airplaneStatus                   = "Airplane mode"
	slowStatus                       = "Slow"
	weakStatus                       = "Weak connection"
)

const (
	updateIntervalSeconds = 5 // interval to update network information, in seconds
)

var (
	connectionIcons  = []rune{'\uf92e', '\uf91e', '\uf921', '\uf924', '\uf927'}
	packetLossIcons  = []rune{'\uf92a', '\uf91f', '\uf922', '\uf925', '\uf928'}
	disconnectedIcon = '\uf92e'
	disabledIcon     = '\uf92d'
)

var lineReturnRegex = regexp.MustCompile(`\r?\n`)

// func getWifi() {
// 	// TODO error checking

// 	output, err := exec.Command("bash", "-c", statusCmd).Output()
// 	if err != nil {
// 		println("Error getting connection status")
// 		return
// 	}

// 	strOutput := string(output)

// 	interfaces := lineReturnRegex.Split(strOutput, -1)
// 	var stringToUse string
// 	for _, v := range interfaces {
// 		if strings.Contains(v, "wifi") {
// 			stringToUse = v
// 			break
// 		}
// 	}

// 	if stringToUse == "" {
// 		return
// 	}

// 	if enabled := strings.Contains(stringToUse, "connect"); enabled {
// 		// ssid := strings.Split(stringToUse, ":")[2]
// 		if strings.Contains(stringToUse, "connecting") {
// 			// see if we're connecting
// 			// block.Set(format.UrgencyLow, disconnectedIcon, "Connecting to "+ssid, "")
// 			return
// 		} else if !strings.Contains(stringToUse, "disconnected") { // make sure we're not disconnected from wifi
// 			// ssid := strings.Split(stringToUse, ":")[2]
// 			// signal, err := getSignalStrength(block)

// 			if err != nil {
// 				// block.Set(format.UrgencyLow, disconnectedIcon, ssid, "")
// 				return
// 			}

// 			// determine which icons we'll use based on
// 			// packetLoss
// 			var icons []rune
// 			if packetLoss() {
// 				icons = packetLossIcons
// 			} else {
// 				icons = connectionIcons
// 			}

// 			// get the icon
// 			iconIndex := len(icons) * signal / 100

// 			// constrains index
// 			if iconIndex >= len(icons) {
// 				iconIndex = len(icons) - 1
// 			}

// 			// block.Set(format.UrgencyNormal, icons[iconIndex], ssid, "")
// 			return
// 		}
// 		// if none of the above, we're disconnected
// 		// block.Set(format.UrgencyLow, disconnectedIcon, "", "")
// 		return
// 	}
// 	// disabled icon
// 	// block.Set(format.UrgencyLow, disabledIcon, "", "")
// }

// func getSignalStrength(block *format.ClassicBlock) (signal int, err error) {
// 	signalOuput, err := exec.Command("bash", "-c", signalCmd).Output()
// 	if err != nil {
// 		// block.Urgency = format.UrgencyLow
// 		block.Icon = disconnectedIcon
// 		block.SetPrimaryText("")
// 		return
// 	}

// 	// get ssid and signal strength
// 	signalStr := lineReturnRegex.ReplaceAllString(strings.Split(string(signalOuput), ":")[1], "")
// 	signal, err = strconv.Atoi(signalStr)
// 	return
// }

// func getEthernet() string {
// 	return ""
// }

func packetLoss(wirelessInterface string) bool {
	cmd := exec.Command("ping", "-c", "2", "-W", "2", "-I", wirelessInterface, "8.8.8.8")
	err := cmd.Run()
	return err != nil
}
