package main

import (
	"github.com/ja7ad/meilisearch-mcp/internal/cli"
)

func main() {
	c := cli.New()
	c.Execute()
}
