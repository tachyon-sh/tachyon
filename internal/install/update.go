package install

import (
	"fmt"
	"os"
	"path/filepath"
	"tachyon/internal/archive"
	"tachyon/internal/env"
)

func UpdatePackage(pkgPath string) error {
	fmt.Println("🔄 Обновление пакета:", pkgPath)

	sitePackages, err := env.GetSitePackagesPath()
	if err != nil {
		return err
	}

	packageName := filepath.Base(pkgPath)
	packageName = packageName[:len(packageName)-len(filepath.Ext(pkgPath))]

	packagePath := filepath.Join(sitePackages, packageName)
	distInfoPath := filepath.Join(sitePackages, packageName+".dist-info")

	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return fmt.Errorf("❌ Ошибка: пакет %s не установлен", packageName)
	}

	tempDistInfo := filepath.Join(os.TempDir(), packageName+".dist-info")
	if _, err := os.Stat(distInfoPath); err == nil {
		err = os.Rename(distInfoPath, tempDistInfo)
		if err != nil {
			return fmt.Errorf("❌ Ошибка сохранения package.dist-info: %v", err)
		}
	}

	err = os.RemoveAll(packagePath)
	if err != nil {
		return fmt.Errorf("❌ Ошибка удаления старой версии: %v", err)
	}

	err = archive.ExtractTPK(pkgPath, sitePackages)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tempDistInfo); err == nil {
		err = os.Rename(tempDistInfo, distInfoPath)
		if err != nil {
			fmt.Println("⚠️ Предупреждение: package.dist-info не удалось восстановить.")
		}
	}

	fmt.Println("✅ Пакет успешно обновлён:", packageName)
	return nil
}