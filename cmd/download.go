package main

import (
	"flag"
	"fmt"
	"sync"

	gv "../govids"
)

var (
	flagFile         string
	flagOutput       string
	flagMaxDownloads int
)

func init() {
	flag.StringVar(&flagFile, "file", "vids.json", "JSON file of all gophervids.com file")
	flag.StringVar(&flagOutput, "output", "output", "Directory to store downloaded videos")
	flag.IntVar(&flagMaxDownloads, "max", 5, "Maximum concurrent downloads to fetch")
	flag.Parse()
}

func main() {
	// Validates input json and output directory exist
	in := gv.ValidatePath(flagFile)
	out := gv.ValidatePath(flagOutput)

	// Read JSON and build video listing
	j := gv.ReadJSON(in)
	videos := gv.NewVideos(j)

	fmt.Println("Downloading", len(videos), "video(s)")
	fmt.Println("------------------------")

	// Used to process downloads concurrently
	var wg sync.WaitGroup
	ch := make(chan struct{}, flagMaxDownloads)
	wg.Add(len(videos))

	// Process downloads
	for _, v := range videos {
		go func(v gv.Video) {
			ch <- struct{}{}
			v.Download(&wg, out)
			<-ch
		}(v)
	}

	// Clean up
	wg.Wait()
	close(ch)

	fmt.Println("------------------------")
	fmt.Println("Finished downloading", len(videos), "video(s)")
}
