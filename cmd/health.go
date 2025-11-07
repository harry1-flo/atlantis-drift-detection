package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the system",
	Long:  `Check the health of the system`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("health check")
	},
}
