package dl

import (
	"crypto/md5"
	"encoding/hex"
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
	var (
		md5List []string
	)

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
		x, _ := os.Open(db.DbPath + file.Name())
		defer x.Close()

		hash := md5.New()
		_, _ = io.Copy(hash, x)
		sum := hex.EncodeToString(hash.Sum(nil)[:])

		for _, md5 := range md5List {
			if md5 == sum {
				os.Remove(db.DbPath + info.Name())
				break
			}
		}
		md5List = append(md5List, sum)

		if dt.Format("2006") != info.ModTime().Format("2006") {
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
