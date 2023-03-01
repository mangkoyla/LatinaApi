package main

import (
	"fmt"
	"os"
	"time"

	"github.com/LalatinaHub/LatinaApi/api/router"
	"github.com/LalatinaHub/LatinaApi/internal/helper"
	latinasub "github.com/LalatinaHub/LatinaSub-go"
	"github.com/LalatinaHub/LatinaSub-go/db"
	"github.com/go-co-op/gocron"
)

func cronJob() {
	schedule := gocron.NewScheduler(time.UTC)

	schedule.Every(6).Hour().Do(func() {
		fmt.Println("Scraping accounts ...")
		helper.LogFuncToFile(func() {
			latinasub.Start()
		}, "scrape.log")
	})

	schedule.StartAsync()
}

func checkDir() {
	_, err := os.Stat("resources")
	if err != nil {
		panic("Could not find resources folder, exiting...")
	}
}

func main() {
	// Check directory
	checkDir()
	db.Init()

	// Set cron job to Download database
	cronJob()

	// Start the router
	router.Start()
}
