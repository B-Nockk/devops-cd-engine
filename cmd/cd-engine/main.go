// cmd/main.go
package main

import (
	"log"

	"cd-engine/internal/ports/in/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
