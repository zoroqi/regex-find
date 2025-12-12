package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/zoroqi/regex-find/internal/app"
)

func main() {
	filePath := flag.String("file", "", "Path to a file to load.")
	flag.Parse()

	var initialText string
	var err error

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			log.Fatalf("Error reading from stdin: %v", readErr)
		}
		initialText = string(bytes)
	} else if *filePath != "" {
		bytes, readErr := os.ReadFile(*filePath)
		if readErr != nil {
			log.Fatalf("Error reading file %s: %v", *filePath, readErr)
		}
		initialText = string(bytes)
	}

	if err = app.New(initialText).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
