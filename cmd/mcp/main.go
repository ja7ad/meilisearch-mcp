package main

import (
	"log"

	"github.com/ja7ad/meilisearch-mcp/internal/cli"
)

func main() {
	c := cli.New()
	if err := c.Execute(); err != nil {
		log.Fatal(err)
	}
}
