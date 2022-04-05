// Package cmd /*
package cmd

import (
	"fmt"
	"telegrambot/internal/care/vars"
	"time"

	"github.com/spf13/cobra"
)

// careCmd represents the care command
var careCmd = &cobra.Command{
	Use:   "care",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("care called")
	},
}

func init() {
	rootCmd.AddCommand(careCmd)
	for {
		time.Sleep(vars.TimeToSleep)
	}
}
