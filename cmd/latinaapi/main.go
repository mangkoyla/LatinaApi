package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mangkoyla/LatinaApi/api/router"
	"github.com/mangkoyla/LatinaApi/common/account"
	"github.com/mangkoyla/LatinaApi/common/account/converter"
	"github.com/mangkoyla/LatinaApi/common/helper"
	latinabot "github.com/mangkoyla/LatinaBot"
	latinasub "github.com/mangkoyla/LatinaSub-go"
	"github.com/go-co-op/gocron"
)

var (
	botToken      = os.Getenv("BOT_TOKEN")
	chatID        = os.Getenv("CHAT_ID")
	sampleTopicID = os.Getenv("SAMPLE_TOPIC_ID")
)

func cronJob() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	schedule := gocron.NewScheduler(loc)
	// schedule.SetMaxConcurrentJobs(1, gocron.RescheduleMode)

	schedule.Cron("30 * * * *").Tag("filter").Do(func() {
		nodes := strings.Split(converter.ToRaw(account.Get("")), "\n")
		if len(nodes) > 500 {
			fmt.Println("Filtering accounts ...")
			latinasub.Start(nodes)
		}
	})

	schedule.Cron("0 */12 * * *").Tag("scrape").Do(func() {
		fmt.Println("Scraping accounts ...")
		helper.LogFuncToFile(func() {
			latinasub.Start([]string{})
		}, "scrape.log")
	})

	// Telegram bot
	if botToken != "" {
		fmt.Println("Starting telegram bot ...")
		go latinabot.Start()

		if chatID != "" {
			var (
				intChatID, _        = strconv.Atoi(chatID)
				intSampleTopicID, _ = strconv.Atoi(sampleTopicID)
			)

			if intSampleTopicID > 0 {
				schedule.Every(3).Hour().Do(func() {
					log.Println("Send VPN sample to channel ...")
					go latinabot.SendVPNToTopic(int64(intChatID), intSampleTopicID)
				})
			}
		}
	}

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

	// Set cron job
	cronJob()

	// Start server
	router.Start()
}
