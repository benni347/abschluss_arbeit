package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func getTemperature(lat, lon string) (float64, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("missing OPENWEATHER_API_KEY environment variable")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	temperature := data["main"].(map[string]interface{})["temp"].(float64)
	return temperature, nil
}

// Greet returns a greeting for the given name
func (a *App) Temprature() string {
	temprature, err := getTemperature("47.4922", "8.7231")
	if err != nil {
		log.Fatal("Error during creation of Temprature: ", err)
	}
	return fmt.Sprintf("The current temprature in Winterthur is %.2f", temprature)
}
