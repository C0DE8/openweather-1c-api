# OpenWeatherMap One Call API w/ golang
A lightweight wrapper for the OpenWeatherMap One Call API

## Usage

First, create an instance of an OpenWeatherMap struct with your APP ID
```go
package main
import owm1c "github.com/C0DE8/openweathermap-1c-api"

owm := owm1c.OpenWeatherMap{API_KEY: os.Getenv("OWM_APP_ID")}
```
