package install

import (
	"fmt"
	"tachyon/internal/archive"
	"tachyon/internal/env"
)

func Package(pkgPath string) error {
	fmt.Println("📦 Установка пакета:", pkgPath)

	sitePackages, err := env.GetSitePackagesPath()
	if err != nil {
		return err
	}
	fmt.Println("📂 Устанавливаем в:", sitePackages)

	err = archive.ExtractTPK(pkgPath, sitePackages)
	if err != nil {
		return err
	}

	fmt.Println("✅ Пакет установлен успешно!") 
	return nil
}