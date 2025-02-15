package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/archive"

	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack <source-dir> -o <output.tpk>",
	Short: "Создаёт .tpk из указанной директории",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ Ошибка: укажите директорию для упаковки")
			os.Exit(1)
		}

		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			fmt.Println("❌ Ошибка: укажите -o <output.tpk>")
			os.Exit(1)
		}

		err := archive.PackTPK(args[0], outputFile)
		if err != nil {
			fmt.Println("❌ Ошибка упаковки:", err)
			os.Exit(1)
		}
	},
}

func init() {
	packCmd.Flags().StringP("output", "o", "", "Файл для сохранения .tpk")
	rootCmd.AddCommand(packCmd)
}