package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"tachyon/internal/install"
	"tachyon/internal/legacy"

	"github.com/spf13/cobra"
)

var deps bool
var channel string

var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Устанавливает пакет",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packagePath := args[0]

		ext := filepath.Ext(packagePath)
		if ext == ".tpk" {
			err := install.Package(packagePath, deps, channel)
			if err != nil {
				fmt.Println("Ошибка установки:", err)
				os.Exit(1)
			}
		} else if ext == ".whl" || ext == ".tar.gz" {
			err := legacy.InstallLegacyPackage(packagePath)
			if err != nil {
				fmt.Println("Ошибка установки легаси-пакета:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("❌ Неподдерживаемый формат:", ext)
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