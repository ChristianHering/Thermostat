package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ChristianHering/Thermostat/utils"
	"github.com/pkg/errors"
)

//WeatherData stores responses from
//OpenWeatherMap's one call endpoint
type WeatherData struct {
	Latitude       float32 `json:"lat"`
	Longitude      float32 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		UnixTime              int64   `json:"dt"`
		SunriseTime           int64   `json:"sunrise"`
		SunsetTime            int64   `json:"sunset"`
		Temperature           float32 `json:"temp"`
		TemperaturePerception float32 `json:"feels_like"`
		AtmosphericPressure   int     `json:"pressure"`
		Humidity              int     `json:"humidity"`
		DewPoint              float32 `json:"dew_point"`
		Cloudiness            int     `json:"clouds"`
		UVIndex               float32 `json:"uvi"`
		Visibility            int     `json:"visibility"`
		WindSpeed             float32 `json:"wind_speed"`
		WindGust              float32 `json:"wind_gust,omitempty"`
		WindDirection         int     `json:"wind_deg"`
		Rain                  struct {
			HourVolume float32 `json:"1h,omitempty"`
		} `json:"rain,omitempty"`
		Snow struct {
			HourVolume float32 `json:"1h,omitempty"`
		} `json:"snow,omitempty"`
		Weather []struct {
			ID          int    `json:"id"`
			Name        string `json:"main"`
			Description string `json:"description"`
			IconID      string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Minutely []struct {
		UnixTime            int64   `json:"dt"`
		PrecipitationVolume float32 `json:"precipitation"`
	} `json:"minutely,omitempty"`
	Hourly []struct {
		UnixTime                  int64   `json:"dt"`
		Temperature               float32 `json:"temp"`
		TemperaturePerception     float32 `json:"feels_like"`
		AtmosphericPressure       int     `json:"pressure"`
		Humidity                  int     `json:"humidity"`
		DewPoint                  float32 `json:"dew_point"`
		Cloudiness                int     `json:"clouds"`
		Visibility                int     `json:"visibility"`
		WindSpeed                 float32 `json:"wind_speed"`
		WindGust                  float32 `json:"wind_gust,omitempty"`
		WindDirection             int     `json:"wind_deg"`
		PrecipitationProbablility float32 `json:"pop"`
		Rain                      struct {
			HourVolume float32 `json:"1h,omitempty"`
		} `json:"rain,omitempty"`
		Snow struct {
			HourVolume float32 `json:"1h,omitempty"`
		} `json:"snow,omitempty"`
		Weather []struct {
			ID          int    `json:"id"`
			Name        string `json:"main"`
			Description string `json:"description"`
			IconID      string `json:"icon"`
		} `json:"weather"`
	} `json:"hourly"`
	Daily []struct {
		UnixTime    int64 `json:"dt"`
		SunriseTime int64 `json:"sunrise"`
		SunsetTime  int64 `json:"sunset"`
		Temperature struct {
			DayTemperature     float32 `json:"day"`
			MinTemperature     float32 `json:"min"`
			MaxTemperature     float32 `json:"max"`
			NightTemperature   float32 `json:"night"`
			EveningTemperature float32 `json:"eve"`
			MorningTemperature float32 `json:"morn"`
		} `json:"temp"`
		TemperaturePerception struct {
			DayTemperature     float32 `json:"day"`
			NightTemperature   float32 `json:"night"`
			EveningTemperature float32 `json:"eve"`
			MorningTemperature float32 `json:"morn"`
		} `json:"feels_like"`
		AtmosphericPressure int     `json:"pressure"`
		Humidity            int     `json:"humidity"`
		DewPoint            float32 `json:"dew_point"`
		WindSpeed           float32 `json:"wind_speed"`
		WindGust            float32 `json:"wind_gust,omitempty"`
		WindDirection       int     `json:"wind_deg"`
		Clouds              int     `json:"clouds"`
		Pop                 float32 `json:"pop"`
		Rain                float32 `json:"rain,omitempty"`
		Snow                float32 `json:"snow,omitempty"`
		UVIndex             float32 `json:"uvi"`
		Weather             []struct {
			ID          int    `json:"id"`
			Name        string `json:"main"`
			Description string `json:"description"`
			IconID      string `json:"icon"`
		} `json:"weather"`
	} `json:"daily"`
}

//getWeatherData returns a struct of weather
//data from OpenWeatherMap's one call API
func getWeatherData() (WD WeatherData, err error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/onecall?lat=" + utils.Config.Latitude + "&lon=" + utils.Config.Longitude + "&units=" + utils.Config.DataUnits + "&appid=" + utils.Config.APIKey)
	if err != nil {
		return WD, errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WD, errors.WithStack(err)
	}

	err = json.Unmarshal(body, &WD)
	if err != nil {
		return WD, errors.WithStack(err)
	}

	return WD, nil
}
