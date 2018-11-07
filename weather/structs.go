package weather

type locationResponse struct {
	Location location `json:"location"`
	Accuracy float32 `json:"accuracy"`
}

type location struct {
	Latitude float32	`json:"lat"`
	Longitude float32	`json:"lng"`
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

