package dl

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func DownloadResource() {
	fmt.Print("Updating Database...")

	dt := time.Now()

	if _, err := os.Stat(db.DbPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(db.DbPath, os.ModePerm)
		} else {
			log.Panic(err)
		}
	}

	out, err := os.Create(db.DbPath + dt.Format("20060102150405") + ".sqlite")
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	resp, err := http.Get(db.DbUrl)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	// Remove unused databases
	files, _ := ioutil.ReadDir("resources/databases/")
	for _, file := range files {
		info, _ := os.Stat(db.DbPath + file.Name())

		if dt.Format("15") != info.ModTime().Format("15") {
			os.Remove(db.DbPath + info.Name())
		}
	}

	fmt.Println("done")
}

// On test
func Scrape() {
	fmt.Println("Starting scraper !")

	scrape := exec.Command("npm", "run", "bg")

	scrape.Dir = "LatinaSub"
	_ = scrape.Run()
}
