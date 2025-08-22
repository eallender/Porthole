package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log/slog"
)

var generateDocsCmd = &cobra.Command{
	Use:    "generate-docs",
	Hidden: true,
	Short:  "Generate documentation for all commands",
	Long:   `Generates markdown documentation for all CLI commands and saves them to the docs/ directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Generating documentation...")

		err := doc.GenMarkdownTree(rootCmd, "./docs/commands/")
		if err != nil {
			slog.Error("Failed to generate markdown docs", "error", err)
			return
		}

		err = doc.GenManTree(rootCmd, &doc.GenManHeader{
			Title:   "PORTHOLE",
			Section: "1",
		}, "./docs/man/")
		if err != nil {
			slog.Error("Failed to generate man pages", "error", err)
			return
		}

		slog.Info("Documentation generated successfully")
		slog.Info("Markdown docs: ./docs/commands/")
		slog.Info("Man pages: ./docs/man/")
	},
}

func init() {
	rootCmd.AddCommand(generateDocsCmd)
}
