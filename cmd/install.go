package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/install"

	"github.com/spf13/cobra"
)

var deps bool

var installCmd = &cobra.Command{
	Use:   "install <package.tpk>",
	Short: "Устанавливает пакет",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packagePath := args[0]

		err := install.Package(packagePath, deps)
		if err != nil {
			fmt.Println("Ошибка установки:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Пакет установлен успешно!")
	},
}

func init() {
	installCmd.Flags().BoolVarP(&deps, "deps", "d", false, "Установить зависимости")
	rootCmd.AddCommand(installCmd)
}