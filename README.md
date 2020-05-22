# openweather-1c-api w/ golang
A lightweight wrapper for the OpenWeather One Call API

## Usage

First, create an instance of an OpenWeatherMap struct with your APP ID
```go
package main
import "github.com/C0DE8/openweather-1c-api"

owm := openweathermap.OpenWeatherMap{API_KEY: os.Getenv("OWM_APP_ID")}
```
