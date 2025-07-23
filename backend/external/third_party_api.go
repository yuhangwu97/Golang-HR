package external

import (
	"encoding/json"
	"time"
)

type ThirdPartyAPI struct {
	client *HTTPClient
}

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

func NewThirdPartyAPI(baseURL string) *ThirdPartyAPI {
	return &ThirdPartyAPI{
		client: NewHTTPClient(baseURL, 30*time.Second),
	}
}

func (api *ThirdPartyAPI) GetWeather(city string) (*WeatherResponse, error) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	response, err := api.client.Get("/weather?city="+city, headers)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(response, &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

func (api *ThirdPartyAPI) GetGeoLocation(ip string) (*GeoLocation, error) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	response, err := api.client.Get("/geo?ip="+ip, headers)
	if err != nil {
		return nil, err
	}

	var location GeoLocation
	if err := json.Unmarshal(response, &location); err != nil {
		return nil, err
	}

	return &location, nil
}