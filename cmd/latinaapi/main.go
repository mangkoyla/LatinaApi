package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaApi/api/router"
	"github.com/LalatinaHub/LatinaApi/common/account"
	"github.com/LalatinaHub/LatinaApi/common/account/converter"
	"github.com/LalatinaHub/LatinaApi/common/helper"
	latinabot "github.com/LalatinaHub/LatinaBot"
	latinasub "github.com/LalatinaHub/LatinaSub-go"
	"github.com/LalatinaHub/LatinaSub-go/db"
	"github.com/go-co-op/gocron"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
	chatID   = os.Getenv("CHAT_ID")
	topicID  = os.Getenv("TOPIC_ID")
)

func cronJob() {
	schedule := gocron.NewScheduler(time.Now().Location())
	schedule.SetMaxConcurrentJobs(1, gocron.RescheduleMode)

	schedule.Cron("30 * * * *").Tag("filter").Do(func() {
		fmt.Println("Filtering accounts ...")
		helper.LogFuncToFile(func() {
			nodes := strings.Split(converter.ToRaw(account.Get("")), "\n")

			if len(nodes) > 500 {
				latinasub.Start(nodes)
			}
		}, "scrape.log")
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

		if chatID != "" && topicID != "" {
			var (
				intChatID, _  = strconv.Atoi(chatID)
				intTopicID, _ = strconv.Atoi(topicID)
			)

			schedule.Every(3).Hour().Do(func() {
				log.Println("Send VPN sample to channel ...")
				go latinabot.SendVPNToTopic(int64(intChatID), intTopicID)
			})
		}
	}

	schedule.StartAsync()
	schedule.RunByTag("scrape")
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

	// Start server
	router.Start()
}
