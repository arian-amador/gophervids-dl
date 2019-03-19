package govids

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/rylio/ytdl"
)

// Video describes structure for a video entry
type Video struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Title string `json:"title"`
}

// NewVideos returns listing of video structs
func NewVideos(b []byte) []Video {
	var videos []Video
	if err := json.Unmarshal(b, &videos); err != nil {
		log.Fatal("Invalid JSON: ", err)
	}

	return videos
}

// Download uses ytdl to download and save a video
func (v *Video) Download(wg *sync.WaitGroup, dir string) {
	defer wg.Done()

	title := v.Filename()
	fullPath := v.FullPath(dir)
	url := v.URL()

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Exists | %s | %s\n", url, title)
		return
	}

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		fmt.Printf("Failed to get video info for %s\n", url)
		return
	}

	file, _ := os.Create(fullPath)
	defer file.Close()

	fmt.Printf("Fetching | %s | %s \n", url, title)
	vid.Download(vid.Formats[0], file)

	return
}

// URL return the youtube url
func (v *Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}

// Filename returns a sanitized title used for the output filename
func (v *Video) Filename() string {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	t := reg.ReplaceAllString(v.Title, "")
	t = strings.ToLower(t)
	t = strings.Replace(t, " ", "-", -1)

	return t
}

// FullPath returns a full output path to save the video
func (v *Video) FullPath(dir string) string {
	return fmt.Sprintf("%s/%s-%s.mp4", dir, v.Date, v.Filename())
}
