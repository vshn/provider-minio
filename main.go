package main

import (
	"fmt"
	"time"
)

var (
	// these variables are populated by Goreleaser when releasing
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")
)

func main() {
	fmt.Printf("version %s, commit %s, date %s\n", version, commit, date)
}
