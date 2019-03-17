package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/rylio/ytdl"
)

var (
	flagFile   string
	flagOutput string
	inFile     string
	outDir     string
)

type video struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Title string `json:"title"`
}

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.Parse()
}

func validatePath(p string) string {
	v, err := filepath.Abs(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(v); os.IsNotExist(err) {
		log.Fatal(err)
	}

	return v
}

func main() {
	inFile = validatePath(flagFile)
	outDir = validatePath(flagOutput)

	jsonFile, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var videos []video
	json.Unmarshal(byteValue, &videos)

	ch := make(chan video, 5)

	var wg sync.WaitGroup

	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				v.download()
			}
		}()
	}

	for _, v := range videos {
		ch <- v
	}

	close(ch)
	wg.Wait()
}

func (v *video) download() {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}
	outFile := reg.ReplaceAllString(v.Title, "")
	outFile = strings.ToLower(outFile)
	outFile = strings.Replace(outFile, " ", "-", -1)
	outFile = fmt.Sprintf("%s/%s-%s.mp4", outDir, v.Date, outFile)

	if _, err := os.Stat(outFile); err == nil {
		fmt.Println("Already downloaded", outFile)
		return
	}

	fmt.Println("Downlaoding", outFile)

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		fmt.Println("Failed to get video info for", outFile)
		return
	}

	file, _ := os.Create(outFile)
	defer file.Close()
	vid.Download(vid.Formats[0], file)

	return
}
