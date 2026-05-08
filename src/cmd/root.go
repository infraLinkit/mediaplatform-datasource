package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long:  `Cobra is a CLI library for Go that empowers applications.`,
	}
)

func init() {

	/**
	 * WEBSERVER SERVICE
	 */
	rootCmd.AddCommand(serverCmd)

	/**
	 * DB MIGRATION (standalone)
	 */
	rootCmd.AddCommand(migrateCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
