package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert is a format conversion tool",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
