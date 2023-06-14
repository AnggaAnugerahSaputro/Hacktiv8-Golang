package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"assignment_3/models"
)

func GetStatusData() *models.Data {
	// Baca data dari file JSON
	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("error while reading data.json file:", err.Error())
	}

	var data models.Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal("error while unmarshalling JSON data:", err.Error())
	}

	return &data
}

func UpdateData() {
	for {
		var data = models.Data{Status: models.Status{}}
		waterMin := 1
		waterMax := 100

		data.Status.Water = rand.Intn(waterMax-waterMin+1) + waterMin
		data.Status.Wind = rand.Intn(waterMax-waterMin+1) + waterMin

		if data.Status.Water < 5 {
			data.WaterStatus = "aman"
		} else if data.Status.Water > 8 {
			data.WaterStatus = "bahaya"
		} else {
			data.WaterStatus = "siaga"
		}

		if data.Status.Wind < 6 {
			data.WindStatus = "aman"
		} else if data.Status.Wind > 15 {
			data.WindStatus = "bahaya"
		} else {
			data.WindStatus = "siaga"
		}

		b, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Fatal("error while marshalling JSON data:", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)
		if err != nil {
			log.Fatal("error while writing data to data.json file:", err.Error())
		}

		time.Sleep(time.Second * 15)
	}
}
