package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/install"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <package>",
	Short: "Удаляет установленный пакет",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]

		err := install.RemovePackage(packageName)
		if err != nil {
			fmt.Println("Ошибка удаления:", err)
			os.Exit(1)
		}

		fmt.Println("❌ Пакет", packageName, "успешно удалён.")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}