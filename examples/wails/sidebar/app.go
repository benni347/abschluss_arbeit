package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

func (a *App) IconGenerator(amount string) ([]string, error) {
	// https://robohash.org/$randomNumber?gravatar=yes&size=500x500
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return nil, err
	}
	icons := make([]string, amountInt)
	for i := 0; i < amountInt; i++ {
		randomNumber := rand.Int()
		filePath := "./frontend/src/assets/images/"
		fileName := fmt.Sprintf("%d.png", randomNumber)
		url := fmt.Sprintf("https://robohash.org/%d?gravatar=yes&size=500x500", randomNumber)
		fullPath := filePath + fileName
		fmt.Printf("Full Path: %s\n", fullPath)
		file, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		_, writeErr := file.Write(body)
		if writeErr != nil {
			return nil, writeErr
		}
		icons[i] = fileName
	}
	return icons, nil
}

// https://opentdb.com/api.php?amount=10&difficulty=hard
