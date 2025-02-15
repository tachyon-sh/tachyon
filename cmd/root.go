package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tachyon",
	Short: "Tachyon - коллайдер",
    Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Используйте `tachyon --help` для списка команд.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
}