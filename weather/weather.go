package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"muse-status/format"
	"net/http"
	"strconv"
	"time"
	"unicode"
)

const (
	updateIntervalMinutes = 3                                  // interval after which to update weather, in minutes
	apiKey                = "d179cc80ed41e8080f9e86356b604ee3" // OpenWeatherMap API key
	units                 = "imperial"
	locationServicesURL   = "https://location.services.mozilla.com/v1/geolocate?key=geoclue"
	openWeatherMapURL     = "http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=%s"
    defaultIcon           = '\uf50f'
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
)

// StartWeatherBroadcast returns a string channel that is fed weather
// information
func StartWeatherBroadcast() chan string {
	channel := make(chan string)

	go func() {
		for {
			loc, err := getLocationJSON()
			if err != nil {
				println("Weather couldn't get the location. Retrying in 10 seconds.")
				time.Sleep(time.Second * 10)
				continue
			}

			report, err := getFullWeatherReport(loc)
			if err != nil || len(report.Weather) <= 0 {
				println("Weather couldn't get a weather report. Retrying in 10 seconds.")
				time.Sleep(time.Second * 10)
				continue
			}

			icon := getWeatherIcon(report)
			weatherString := getWeatherString(report)

			channel <- icon + "  " + weatherString

			time.Sleep(time.Minute * updateIntervalMinutes)
		}
	}()

	return channel
}

func getWeatherIcon(report fullWeatherReport) string {
    if len(report.Weather) <= 0 {
        return ""
    }

    iconString := report.Weather[0].Icon
    if icon, ok := weatherIcons[iconString]; ok {
        return string(icon)
    }
	return iconString
}

func getWeatherString(report fullWeatherReport) string {
    if len(report.Weather) <= 0 {
        return ""
    }

	// basically round degrees to the nearest int and add the degree sign
	degrees := strconv.Itoa(int(report.Main.Temp+0.5)) + "°"

	// capitalize the first letter in the description
	desc := []rune(report.Weather[0].Description)
	desc[0] = unicode.ToUpper(desc[0])

	return degrees + " " + format.Dim(string(desc))
}

func getLocationJSON() (location, error) {
	res, err := http.Get(locationServicesURL)
	if err != nil {
		return location{}, err
	}

	defer res.Body.Close()

	// get response as a []byte
	var locResponse locationResponse
	resBodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return location{}, err
	}

	// decode json
	json.Unmarshal(resBodyStr, &locResponse)

	return locResponse.Location, nil
}

func getFullWeatherReport(loc location) (report fullWeatherReport, err error) {
	reqURL := fmt.Sprintf(openWeatherMapURL, loc.Latitude, loc.Longitude, apiKey, units)
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
