package dl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func DownloadResource() {
	fmt.Print("Updating Database...")
	_, err := os.Stat(db.DbPath)
	if err != nil {
		os.Remove(db.DbPath)
	}

	out, err := os.Create(db.DbPath)
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

	fmt.Println("done")
}

// On test
func Scrape() {
	fmt.Println("Starting scraper !")

	scrape := exec.Command("npm", "run", "bg")

	scrape.Dir = "LatinaSub"
	_ = scrape.Run()
}
