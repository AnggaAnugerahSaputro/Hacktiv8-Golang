package models

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status       `json:"status"`
	WaterStatus  string `json:"water_status"`
	WindStatus   string `json:"wind_status"`
}
