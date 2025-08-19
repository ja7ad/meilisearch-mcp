package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CLI struct {
	root *cobra.Command
}

func New() *CLI {
	c := &CLI{
		root: &cobra.Command{
			Use:   "meilisearch-mcp",
			Short: "Meilisearch MCP tool",
		},
	}

	debug := c.root.PersistentFlags().Bool("debug", false, "Enable debug mode")

	c.serve(*debug)

	return c
}

func (c *CLI) Root() *cobra.Command { return c.root }

func (c *CLI) Execute() error {
	if err := c.root.Execute(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
