package weather

type location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type fullWeatherReport struct {
	Sys     sunTimeData      `json:"sys"`
	Weather []weatherDetails `json:"weather"`
	Main    weatherMain      `json:"main"`
	Wind    weatherWind      `json:"wind"`
}

type sunTimeData struct {
	Sunrise int64 `json:"sunrise"`
	Sunset  int64 `json:"sunset"`
}

type weatherDetails struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type weatherWind struct {
	Speed float32 `json:"speed"`
	Deg   float32 `json:"deg"`
}

type weatherMain struct {
	Temp float32 `json:"temp"`
}
