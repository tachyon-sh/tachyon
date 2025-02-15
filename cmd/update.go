package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/install"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <package.tpk>",
	Short: "Обновляет установленный пакет",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packagePath := args[0]

		err := install.UpdatePackage(packagePath)
		if err != nil {
			fmt.Println("Ошибка обновления:", err)
			os.Exit(1)
		}

		fmt.Println("🔄 Пакет успешно обновлён:", packagePath)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}