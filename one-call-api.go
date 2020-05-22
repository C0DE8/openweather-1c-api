package openweathermap_1c_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Schema   string = "https"
	Domain   string = "api.openweathermap.org"
	PathTmpl string = "/data/%v/"
	Version  string = "2.5"
	Resource string = "onecall?lat=%v&lon=%v&exclude=hourly,daily&appid=%v"
)

var MissingApiKey = errors.New("no API key present")

func DefaultClient(timeout int) Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

// Client : an interface to allow other clients to be injected (and of course for testing purposes)
type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type OpenWeatherMapOneCallAPI struct {
	Client Client
	ApiKey string
	Unit   string
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type TZ struct {
	Timezone string `json:"timezone"`
	Offset   int32  `json:"timezone_offset"`
}

type Sun struct {
	Rise uint32 `json:"sunrise"`
	Set  uint32 `json:"sunset"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
}

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float32 `json:"wind_speed"`
	Deg   uint16  `json:"wind_deg"`
}

type Rain struct {
	Threehr int `json:"3h"`
}

type Uvi struct {
	Value float32 `json:"uvi"`
}

type Section struct { // current, minutely, hourly
	Dt int `json:"dt"`
	Sun
	Main    `json:"main"`
	Wind    `json:"wind"`
	Rain    `json:"rain"`
	Clouds  `json:"clouds"`
	Uvi     `json:"uvi"`
	Weather []Weather `json:"weather"`
}

type Current struct {
	Section
}

type Hourly struct {
	Section
}

type Minutely struct {
	Section
}

type Response struct {
	Coordinate
	TZ
	Current
	Hourly
	Minutely
}

func (owm *OpenWeatherMapOneCallAPI) execApiRequest(url string) ([]byte, error) {
	// build an http client so we can have control over timeout
	client := &http.Client{
		Timeout: time.Duration(owm.ClientTimeout) * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// defer the closing of the response body (with error handling)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	// read the http response body into a byte stream
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (owm *OpenWeatherMapOneCallAPI) GetWeatherFromLatLon(lat, lon float64, time int64) (*Response, error) {
	if owm.ApiKey == "" {
		// No API keys present, return error
		return nil, MissingApiKey
	}

	url := fmt.Sprintf(getBasicUrlTemplate(), Version, Resource, lat, lon, time, owm.ApiKey)

	return owm.processRequest(url)
}

func (owm *OpenWeatherMapOneCallAPI) processRequest(url string) (*Response, error) {
	body, err := owm.execApiRequest(url)
	if err != nil {
		return nil, err
	}

	var cwr Response

	// (try to) unmarshal the (JSON) byte stream into a Go data type
	err = json.Unmarshal(body, &cwr)
	if err != nil {
		return nil, err
	}

	return &cwr, nil
}

func getBasicUrlTemplate() string {
	return Schema + "://" + Domain + PathTmpl + Resource
}

/*
 {
  "lat": 33.44,
  "lon": -94.04,
  "timezone": "America/Chicago",
  "timezone_offset": -18000,
  "current": {
    "dt": 1588935779,
    "sunrise": 1588936856,
    "sunset": 1588986260,
    "temp": 16.75,
    "feels_like": 16.07,
    "pressure": 1009,
    "humidity": 93,
    "dew_point": 15.61,
    "uvi": 8.97,
    "clouds": 90,
    "visibility": 12874,
    "wind_speed": 3.6,
    "wind_deg": 280,
    "weather": [
      {
        "id": 501,
        "main": "Rain",
        "description": "moderate rain",
        "icon": "10n"
      },
      {
        "id": 200,
        "main": "Thunderstorm",
        "description": "thunderstorm with light rain",
        "icon": "11n"
      }
    ],
    "rain": {
      "1h": 2.79
    }
  },
   "minutely": [
    {
      "dt": 1588935780,
      "precipitation": 2.789
    },
    ...
  },
  "hourly": [
      {
      "dt": 1588935600,
      "temp": 16.75,
      "feels_like": 13.93,
      "pressure": 1009,
      "humidity": 93,
      "dew_point": 15.61,
      "clouds": 90,
      "wind_speed": 6.66,
      "wind_deg": 203,
      "weather": [
        {
          "id": 501,
          "main": "Rain",
          "description": "moderate rain",
          "icon": "10n"
        }
      ],
      "rain": {
        "1h": 2.92
      }
    },
    ...
  }
    "daily": [
        {
      "dt": 1588960800,
      "sunrise": 1588936856,
      "sunset": 1588986260,
      "temp": {
        "day": 22.49,
        "min": 10.96,
        "max": 22.49,
        "night": 10.96,
        "eve": 18.45,
        "morn": 18.14
      },
      "feels_like": {
        "day": 18.72,
        "night": 6.53,
        "eve": 16.34,
        "morn": 16.82
      },
      "pressure": 1014,
      "humidity": 60,
      "dew_point": 14.35,
      "wind_speed": 7.36,
      "wind_deg": 342,
      "weather": [
        {
          "id": 502,
          "main": "Rain",
          "description": "heavy intensity rain",
          "icon": "10d"
        }
      ],
      "clouds": 68,
      "rain": 15.38,
      "uvi": 8.97
    },
    ...
    }
*/
