package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)


// WeatherResponse represents the response from Open Weather API
type WeatherResponse struct {
    Weather []struct {
        Main string `json:"main"`
    } `json:"weather"`
    Main struct {
        Temp float64 `json:"temp"`
    } `json:"main"`
}

// WeatherInfo represents the custom weather response
type WeatherInfo struct {
    Condition string `json:"condition"`
    TemperatureStatus string `json:"temperature_status"`
}

func getTemperatureStatus(temp float64) string {
    if temp < 18 {
        return "cold"
    } else if temp >= 18 && temp < 28 {
        return "moderate"
    } else {
        return "hot"
    }
}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    if lat == "" || lon == "" {
        http.Error(w, "lat and lon query parameters are required", http.StatusBadRequest)
        return
    }

    // Call the Open Weather API
    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        log.Fatal("API_KEY environment variable is not set")
    }

    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, apiKey)
    log.Println("get data:", url)
    
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var weatherResp WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
        http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
        return
    }

    tempStatus := getTemperatureStatus(weatherResp.Main.Temp)
    weatherInfo := WeatherInfo{
        Condition: weatherResp.Weather[0].Main,
        TemperatureStatus: tempStatus,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(weatherInfo); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/weather", getWeatherHandler)
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
