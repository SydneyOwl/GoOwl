package cmd

import (
	"errors"
	"fmt"
	"github.com/sydneyowl/GoOwl/cmd/checkenv"
	"github.com/sydneyowl/GoOwl/cmd/run"
	"github.com/sydneyowl/GoOwl/common/global"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "GoOwl",
	Short:   "GoOwl",
	Version: fmt.Sprintf("%s(%s)", global.Version, global.Status),
	Long:    `GoOwl - A simple CI/CD Platform!`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("at least one arg is required!")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func init() {
	rootCmd.AddCommand(checkenv.StartCmd)
	rootCmd.AddCommand(run.StartCmd)
}

// tip displays when input command cannot be identified.
func tip() {
	fmt.Printf("GoOwl Ver:%s (%s).\n", global.Version, global.Status)
	fmt.Println("Use GoOwl -h for help.")
}

// Execute is the entrance of Gowl.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
