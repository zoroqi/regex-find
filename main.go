package main

import (
	"log"

	"github.com/zoroqi/regex-find/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
