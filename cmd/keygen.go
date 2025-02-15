package cmd

import (
	"fmt"
	"os"
	"tachyon/internal/security"

	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Генерирует Ed25519 ключи для подписи пакетов",
	Run: func(cmd *cobra.Command, args []string) {
		err := security.GenerateKeys()
		if err != nil {
			fmt.Println("Ошибка генерации ключей:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}