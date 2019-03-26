package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	gv "github.com/arian-amador/gophervidsdl/gophervids"
	"github.com/arian-amador/gophervidsdl/utils"
)

var (
	flagRemote       bool
	flagFile         string
	flagOutput       string
	flagMaxDownloads int
	flagDebug        bool
)

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.BoolVar(&flagRemote, "remote", false, "Get gophervids.com json listing")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.IntVar(&flagMaxDownloads, "max", 5, "Maximum concurrent downloads to fetch")
	flag.BoolVar(&flagDebug, "debug", false, "Show progress during download process")
	flag.Parse()
}

func download(videos []gv.Video, out string) {
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

	return
}

func main() {
	var videos []gv.Video
	var err error

	out, _ := filepath.Abs(flagOutput)
	if err := utils.ValidatePath(out); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	if flagRemote {
		videos, err = gv.NewRemoteJSON()
	} else {
		videos, err = gv.NewLocalJSON(flagFile)
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	download(videos, out)

	if flagDebug {
		fmt.Println("Finished downloading", len(videos), "video(s)")
	}

	os.Exit(0)
}
