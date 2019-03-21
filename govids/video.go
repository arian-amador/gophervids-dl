package govids

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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
func (v *Video) Download(o string, debug bool) error {
	if debug {
		if _, err := os.Stat(o); err == nil {
			return fmt.Errorf("Exists | %s | %s", v.URL(), v.Title)
		} else {
			fmt.Printf("Fetching | %s | %s \n", v.URL(), v.Title)
		}
	}

	vid, err := ytdl.GetVideoInfo(v.URL())
	if err != nil {
		i := strings.Index(err.Error(), ":")
		err := fmt.Errorf("Error | %s | %s | %s", err.Error()[i+1:], v.URL(), v.Title)
		return err
	}
	if len(vid.Formats) == 0 {
		return fmt.Errorf("Error No videos found at %s", v.URL())
	}

	file, err := os.Create(o)
	if err != nil {
		return err
	}
	defer file.Close()

	vid.Download(vid.Formats[0], file)

	return nil
}

// Author returns the videos channel name
func (v *Video) Author() string {
	url := v.URL()

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return ""
	}

	return vid.Author
}

// URL return the youtube url
func (v *Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}

// Filename returns a sanitized title used for the output filename
func (v *Video) Filename() string {
	return Sanitize(v.Title)
}

// FullPath returns a full output path to save the video
func (v *Video) FullPath(p string) string {
	a := Sanitize(v.Author())

	if a != "" {
		p = p + string(os.PathSeparator) + a
		if err := ValidatePath(p); err != nil {
			os.MkdirAll(p, os.ModePerm)
		}
	}

	p = p + string(os.PathSeparator) + v.Date + "-" + v.Filename()

	return p
}
