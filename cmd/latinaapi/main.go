package main

import (
	"os"
	"time"

	"github.com/LalatinaHub/LatinaApi/api/router"
	"github.com/LalatinaHub/LatinaApi/cmd/dl"
	"github.com/LalatinaHub/LatinaApi/internal/db"
	"github.com/go-co-op/gocron"
	_ "github.com/mattn/go-sqlite3"
)

func cronJob() {
	schedule := gocron.NewScheduler(time.UTC)
	schedule.Every(30).Minutes().Do(func() {
		dl.DownloadResource()
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

	// Set cron job to Download database
	cronJob()

	// Initialize databse
	db.Database.Init()

	// Start the router
	router.Start()
}
