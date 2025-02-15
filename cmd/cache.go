package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/cache"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Управление кешем пакетов",
}

var cacheListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать кешированные пакеты",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := cache.ListCache()
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}

		if len(files) == 0 {
			fmt.Println("🗂 Кеш пуст.")
			return
		}

		fmt.Println("📦 Закешированные пакеты:")
		for _, file := range files {
			fmt.Println("  -", file)
		}
	},
}

func init() {
	cacheCmd.AddCommand(cacheListCmd)
	rootCmd.AddCommand(cacheCmd)
}