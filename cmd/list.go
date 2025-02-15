package cmd

import (
	"fmt"
	"tachyon/internal/env"
	"tachyon/internal/install"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Выводит список установленных пакетов",
	Run: func(cmd *cobra.Command, args []string) {
		sitePackages, err := env.GetSitePackagesPath()
		if err != nil {
			fmt.Println("Ошибка определения среды:", err)
			return
		}

		packages, err := install.ListPackages(sitePackages)
		if err != nil {
			fmt.Println("Ошибка при получении списка пакетов:", err)
			return
		}

		fmt.Println("Установленные пакеты:")
		for _, pkg := range packages {
			fmt.Println(" -", pkg)
		}
	},
}