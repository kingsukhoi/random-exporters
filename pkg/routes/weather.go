package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"random-exporters/pkg/formatter/OpenMetrics"
	"random-exporters/pkg/models"
	"strconv"
	"strings"
)

const baseUrl = "https://api.weatherapi.com/v1"

func OpenMetricsWeatherCurrent(c *gin.Context) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		err := errors.New("WEATHER_API_KEY env var missing")
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	loc := c.Query("loc")
	if loc == "" {
		err := errors.New("loc query param missing")
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=yes", baseUrl, apiKey, loc)

	client := http.Client{}
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", url, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//req.Header.Add("Accept-Encoding", "gzip")
	res, err := client.Do(req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	currentWeather := models.WeatherApiCurrent{}

	err = json.Unmarshal(resBytes, &currentWeather)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tempc := OpenMetrics.Metric{
		Name: "temperature_feels_like_c",
		Type: OpenMetrics.Gauge,
		Labels: map[string]string{
			"city":    currentWeather.Location.Name,
			"country": currentWeather.Location.Country,
		},
		Value: strconv.FormatFloat(currentWeather.Current.FeelslikeC, 'g', -1, 64),
		Help:  "Current Temperature in Celcius",
	}
	pm10 := OpenMetrics.Metric{
		Name: "PM10",
		Type: OpenMetrics.Gauge,
		Labels: map[string]string{
			"city":    currentWeather.Location.Name,
			"country": currentWeather.Location.Country,
		},
		Value: strconv.FormatFloat(currentWeather.Current.AirQuality.Pm10, 'g', -1, 64),
		Help:  "Current PM10 measurement",
	}
	pm25 := OpenMetrics.Metric{
		Name: "PM25",
		Type: OpenMetrics.Gauge,
		Labels: map[string]string{
			"city":    currentWeather.Location.Name,
			"country": currentWeather.Location.Country,
		},
		Value: strconv.FormatFloat(currentWeather.Current.AirQuality.Pm25, 'g', -1, 64),
		Help:  "Current PM25 measurement",
	}

	rtnMe := strings.Builder{}
	rtnMe.WriteString(tempc.String())
	rtnMe.WriteString(pm10.String())
	rtnMe.WriteString(pm25.String())

	c.Data(http.StatusOK, "application/openmetrics-text; version=1.0.0; charset=utf-8", []byte(rtnMe.String()))

}
