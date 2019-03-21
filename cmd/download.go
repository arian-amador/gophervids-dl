package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	gv "github.com/arian-amador/gophervidsdl"
)

var (
	flagRemote       bool
	flagFile         string
	flagOutput       string
	flagMaxDownloads int
	flagDebug        bool
)

const (
	remoteURL = "http://gophervids.appspot.com/static/vids.json"
)

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.BoolVar(&flagRemote, "remote", false, "Get gophervids.com json listing")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.IntVar(&flagMaxDownloads, "max", 5, "Maximum concurrent downloads to fetch")
	flag.BoolVar(&flagDebug, "debug", false, "Show progress during download process")
	flag.Parse()
}

func main() {
	if flagRemote {
		if err := download(flagFile, remoteURL); err != nil {
			log.Fatal("Unable to download video json file")
		}
	}

	// Validate the JSON including the video listing exists
	in, _ := filepath.Abs(flagFile)
	if err := gv.ValidatePath(in); err != nil {
		log.Fatal(err)
	}

	// Validate the output directory exists
	out, _ := filepath.Abs(flagOutput)
	if err := gv.ValidatePath(out); err != nil {
		log.Fatal(err)
	}

	// Read JSON and build video listing
	j := gv.ReadJSON(in)
	videos := gv.NewVideos(j)

	// Used to process downloads concurrently
	var wg sync.WaitGroup
	ch := make(chan struct{}, flagMaxDownloads)
	wg.Add(len(videos))

	// Process downloads
	for _, v := range videos {
		go func(v gv.Video) {
			defer wg.Done()

			ch <- struct{}{}
			p := v.FullPath(out)
			if err := v.Download(p, flagDebug); err != nil {
				fmt.Println(err)
			}
			<-ch
		}(v)
	}

	wg.Wait()
	close(ch)

	if flagDebug {
		fmt.Println("Finished downloading", len(videos), "video(s)")
	}
	os.Exit(0)
}

func download(f string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(f)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
