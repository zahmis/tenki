package ps


type TotalInfo struct {
	Name string `json:"name"`
	Weather Weather`json:"weather"`
	Main Main `json:"main"`
	Wind Wind `json:"wind"`
	Clouds Clouds `json:"clouds"`
	Sys Sys `json:"sys"`
}

type Weather struct {
	Main string `json:"main"`
	Description string `json:"description"`
}

type Main struct {
	Temp float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg float64 `json:"deg"`
}

type Clouds struct {
	All float64 `json:"all"`
}

type Sys struct {
	Country string `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset float64 `json:"sunset"`
}