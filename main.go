package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/zoroqi/regex-find/internal/app"
	"golang.design/x/clipboard"
)

const historyEnvVar = "REGEX_FIND_HISTORY_FILE"

func main() {

	filePath := flag.String("file", "", "Path to a file to load.")
	flag.StringVar(filePath, "f", "", "Path to a file to load.")

	useClipboard := flag.Bool("clipboard", false, "Read initial text from the system clipboard.")
	flag.BoolVar(useClipboard, "c", false, "Read initial text from the system clipboard.")

	historyFile := flag.String("history-file", "", fmt.Sprintf("Path to the history file. Overrides the %s environment variable.", historyEnvVar))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A TUI tool for interactively developing and testing regular expressions.\n\n")
		fmt.Fprintf(os.Stderr, "Input can be provided from a file (--file), clipboard (--clipboard), or stdin.\n")
		fmt.Fprintf(os.Stderr, "Example: cat my_text.log | %s\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Determine history file path
	historyPath := *historyFile
	if historyPath == "" {
		historyPath = os.Getenv(historyEnvVar)
	}

	var initialText string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 { // Data from stdin
		bytes, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			log.Fatalf("Error reading from stdin: %v", readErr)
		}
		initialText = string(bytes)
	} else if *useClipboard && *filePath != "" { // Mutually exclusive check
		log.Fatalf("Error: Cannot use both --clipboard and --file flags simultaneously.")
	} else if *useClipboard { // Read from clipboard
		// Initialize clipboard, this might take some time on Wayland
		if initErr := clipboard.Init(); initErr != nil {
			log.Fatalf("Error initializing clipboard: %v", initErr)
		}
		bytes := clipboard.Read(clipboard.FmtText)
		if bytes == nil {
			initialText = "" // No content in clipboard
		} else {
			initialText = string(bytes)
		}
	} else if *filePath != "" { // Read from file
		bytes, readErr := os.ReadFile(*filePath)
		if readErr != nil {
			log.Fatalf("Error reading file %s: %v", *filePath, readErr)
		}
		initialText = string(bytes)
	}

	appInstance, err := app.New(initialText, historyPath)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	if err := appInstance.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	// Save history on exit
	appInstance.UpdateHistoryWithCurrentRegex()
	if err := appInstance.SaveHistory(); err != nil {
		log.Printf("Warning: could not save history: %v", err)
	}

	fmt.Println(appInstance.GetRegexInput())
}
