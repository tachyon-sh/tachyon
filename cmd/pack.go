package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/archive"

	"github.com/spf13/cobra"
)

var outputPath string

var packCmd = &cobra.Command{
	Use:   "pack <folder>",
	Short: "Упаковывает директорию в .tpk",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFolder := args[0]

		if outputPath == "" {
			outputPath = inputFolder + ".tpk"
		}

		err := archive.CreateTPK(inputFolder, outputPath)
		if err != nil {
			fmt.Println("Ошибка упаковки:", err)
			os.Exit(1)
		}

		fmt.Println("Пакет успешно упакован:", outputPath)
	},
}

func init() {
	packCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Выходной файл .tpk")
	rootCmd.AddCommand(packCmd)
}