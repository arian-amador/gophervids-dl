# Gopher Video Downloader [![Build Status](https://api.travis-ci.org/arian-amador/gophervidsdl.svg)](https://travis-ci.org/arian-amador/gophervidsdl)

Go utility to download and store all of the videos listed on http://gophervids.appspot.com

## Usage

```bash
go get -t github.com/arian-amador/gophervidsdl
go build cmd/download.go
```

## Flags

```bash
  -debug
        Show progress during download process
  -file string
        JSON file of all gophervids.com file (default "vids.json")
  -max int
        Maximum concurrent downloads to fetch (default 5)
  -output string
        Directory to store downloaded videos (default "output")
  -remote
        Get gophervids.com json listing
```
