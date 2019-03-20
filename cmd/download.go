package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	gv "../govids"
)

var (
	flagFile         string
	flagOutput       string
	flagMaxDownloads int
	flagDebug        bool
)

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.IntVar(&flagMaxDownloads, "max", 5, "Maximum concurrent downloads to fetch")
	flag.BoolVar(&flagDebug, "debug", false, "Show progress during download process")
	flag.Parse()
}

func main() {
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
				fmt.Errorf("%s | %s", err, v.Title)
			}
			<-ch
		}(v)
	}

	wg.Wait()
	close(ch)
	os.Exit(0)
}
