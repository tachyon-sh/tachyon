package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/install"

	"github.com/spf13/cobra"
)

var deps bool
var channel string

var installCmd = &cobra.Command{
	Use:   "install <package.tpk>",
	Short: "Устанавливает пакет",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packagePath := args[0]

		err := install.Package(packagePath, deps, channel)
		if err != nil {
			fmt.Println("Ошибка установки:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Пакет установлен успешно!")
	},
}

func init() {
	installCmd.Flags().BoolVarP(&deps, "deps", "d", false, "Установить зависимости")
	installCmd.Flags().StringVarP(&channel, "channel", "c", "stable", "Выбрать канал (stable, beta, nightly)")
	rootCmd.AddCommand(installCmd)
}