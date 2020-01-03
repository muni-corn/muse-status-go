package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "github.com/muni-corn/muse-status/format"
	"io/ioutil"
	"net/http"
	"strconv"
	// "time"
	"unicode"
)

const (
	updateIntervalMinutes     = 10 // interval after which to update weather, in minutes
	units                     = "imperial"
	ipStackURLTemplate        = "http://api.ipstack.com/%s?access_key=%s&format=1" // ip, then key
	ipStackKey                = "9c237911bdacce2e8c9a021d9b4c1317"
	openWeatherMapURLTemplate = "http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=%s"
	openWeatherMapKey         = "d179cc80ed41e8080f9e86356b604ee3" // OpenWeatherMap API key
	defaultIcon               = '\uf50f'
)

var (
	weatherIcons = map[string]rune{
		"01d": '',
		"01n": '',
		"02d": '',
		"02n": '',
		"03d": '',
		"03n": '',
		"04d": '',
		"04n": '',
		"09d": '',
		"09n": '',
		"10d": '',
		"10n": '',
		"11d": '',
		"11n": '',
		"13d": '',
		"13n": '',
		"50d": '',
		"50n": '',
	}
	// nerd font
	// weatherIcons = map[string]rune{
	// 	"01d": '\ue30d',
	// 	"01n": '\ue32b',
	// 	"02d": '\ue30c',
	// 	"02n": '\ue379',
	// 	"03d": '\ue302',
	// 	"03n": '\ue37e',
	// 	"04d": '\ue33d',
	// 	"04n": '\ue33d',
	// 	"09d": '\ue309',
	// 	"09n": '\ue326',
	// 	"10d": '\ue308',
	// 	"10n": '\ue325',
	// 	"11d": '\ue305',
	// 	"11n": '\ue322',
	// 	"13d": '\ue30a',
	// 	"13n": '\ue327',
	// 	"50d": '\ue303',
	// 	"50n": '\ue313',
	// }
)

// StartWeatherBroadcast returns a string channel that is fed weather
// information
// func StartWeatherBroadcast() chan *format.ClassicBlock {
// 	channel := make(chan *format.ClassicBlock)
// 	block := &format.ClassicBlock{Name: "weather"}

// 	go func() {
// 		for {
// 			loc, err := getLocation()
// 			if err != nil {
// 				println("Weather couldn't get the location. Retrying in 10 seconds.")
// 				time.Sleep(time.Second * 10)
// 				continue
// 			}

// 			report, err := getFullWeatherReport(*loc)
// 			if err != nil || len(report.Weather) <= 0 {
// 				println("Weather couldn't get a weather report. Retrying in 10 seconds.")
// 				time.Sleep(time.Second * 10)
// 				continue
// 			}

// 			icon := getWeatherIcon(report)
// 			temperature := getTemperatureString(report)
// 			description := getWeatherDescription(report)

// 			block.Set(format.UrgencyNormal, icon, temperature, description)

// 			channel <- block

// 			time.Sleep(time.Minute * updateIntervalMinutes)
// 		}
// 	}()

// 	return channel
// }

func getExternalIP() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

func getWeatherIcon(report fullWeatherReport) rune {
	if len(report.Weather) <= 0 {
		return ' '
	}

	iconString := report.Weather[0].Icon
	if icon, ok := weatherIcons[iconString]; ok {
		return icon
	}
	return ' '
}

func getTemperatureString(report fullWeatherReport) string {
	if len(report.Weather) <= 0 {
		return ""
	}

	// basically round degrees to the nearest int and add the degree sign
	degrees := strconv.Itoa(int(report.Main.Temp+0.5)) + "°"
	return degrees
}

func getWeatherDescription(report fullWeatherReport) string {
	if len(report.Weather) <= 0 {
		return ""
	}

	// capitalize the first letter in the description
	desc := []rune(report.Weather[0].Description)
	desc[0] = unicode.ToUpper(desc[0])

	return string(desc)
}

func getLocation() (*WeatherLocation, error) {
	ip, err := getExternalIP()
	if err != nil {
		return nil, err
	}
	// println("ip: " + ip)

	url := fmt.Sprintf(ipStackURLTemplate, ip, ipStackKey)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// get response as a []byte
	resBodyStr, err := ioutil.ReadAll(res.Body)
	// println(string(resBodyStr))
	if err != nil {
		return nil, err
	}

	// decode json
	var loc WeatherLocation
	json.Unmarshal(resBodyStr, &loc)

	// println(fmt.Sprintf("location: %f, %f", loc.Latitude, loc.Longitude))

	return &loc, nil
}

func getFullWeatherReport(loc *WeatherLocation) (report fullWeatherReport, err error) {
	if loc == nil {
		return fullWeatherReport{}, fmt.Errorf("nil location given, can't get weather report")
	}
	reqURL := fmt.Sprintf(openWeatherMapURLTemplate, loc.Latitude, loc.Longitude, openWeatherMapKey, units)
	res, err := http.Get(reqURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// convert response body to string
	resBodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(resBodyStr, &report)
	if err != nil {
		return
	}

	return
}
