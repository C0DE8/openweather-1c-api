package main

import (
	"fmt"
	owm1c "openweathermap-1c-api"
)

func main() {
	var owm = owm1c.OpenWeatherMapOneCallAPI{
		Client: owm1c.DefaultClient(5),
		ApiKey: "{YOUR API KEY}",
		Unit:   "metric", // = celsius
	}

	owmResponse, err := owm.GetWeatherFromLatLon(60.99, 30.9, 1586468027)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	fmt.Println(owmResponse)
}
