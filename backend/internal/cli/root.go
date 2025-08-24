package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Release Calendar backend CLI",
}

func Execute() error { return rootCmd.Execute() }

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}
