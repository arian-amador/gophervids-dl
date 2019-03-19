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
	flagFile         string
	flagOutput       string
	flagMaxDownloads int
	flagVerbose      bool
)

type video struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Title string `json:"title"`
}

func (v *video) url() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}

func (v *video) title() string {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	t := reg.ReplaceAllString(v.Title, "")
	t = strings.ToLower(t)
	t = strings.Replace(t, " ", "-", -1)

	return t
}

func (v *video) fullPath(dir string) string {
	return fmt.Sprintf("%s/%s-%s.mp4", dir, v.Date, v.title())
}

func (v *video) download(wg *sync.WaitGroup, dir string) {
	defer wg.Done()

	title := v.title()
	fullPath := v.fullPath(dir)
	url := v.url()

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Exists | %s | %s\n", url, title)
		return
	}

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		fmt.Printf("Failed to get video info for %s\n", url)
		return
	}

	if flagVerbose {
		fmt.Printf("%+v\n", vid)
	}

	file, _ := os.Create(fullPath)
	defer file.Close()

	fmt.Printf("Fetching | %s | %s \n", url, title)
	vid.Download(vid.Formats[0], file)

	return
}

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.IntVar(&flagMaxDownloads, "max", 5, "Maximum concurrent downloads to fetch")
	flag.BoolVar(&flagVerbose, "verbose", false, "Show details of each video as it's being processed")
	flag.Parse()
}

func main() {
	// Validates input json and output directory exist
	in := validatePath(flagFile)
	out := validatePath(flagOutput)

	// Read JSON and build video listing
	videos := readJSON(in)

	fmt.Println("Downloading", len(videos), "video(s)")
	fmt.Println("------------------------")

	// Used to process downloads concurrently
	var wg sync.WaitGroup
	ch := make(chan struct{}, flagMaxDownloads)
	wg.Add(len(videos))

	// Process downloads
	for _, v := range videos {
		go func(v video) {
			ch <- struct{}{}
			v.download(&wg, out)
			<-ch
		}(v)
	}

	// Clean up
	wg.Wait()
	close(ch)

	fmt.Println("------------------------")
	fmt.Println("Finished downloading", len(videos), "video(s)")
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

func readJSON(in string) []video {
	fmt.Println("Reading JSON")

	jsonFile, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var videos []video
	if err := json.Unmarshal(byteValue, &videos); err != nil {
		log.Fatal("Invalid JSON: ", err)
	}

	return videos
}
