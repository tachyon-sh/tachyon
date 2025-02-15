package legacy

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func InstallLegacyPackage(pkgPath string) error {
	ext := filepath.Ext(pkgPath)
	fmt.Println("📦 Установка легаси-пакета:", pkgPath)

	var cmd *exec.Cmd
	if ext == ".whl" {
		cmd = exec.Command("pip", "install", "--no-cache-dir", pkgPath)
	} else if ext == ".tar.gz" {
		cmd = exec.Command("pip", "install", "--no-cache-dir", pkgPath)
	} else {
		return fmt.Errorf("❌ Неподдерживаемый формат: %s", ext)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Ошибка установки:", string(output))
		return err
	}

	fmt.Println("✅ Легаси-пакет установлен!")
	return nil
}