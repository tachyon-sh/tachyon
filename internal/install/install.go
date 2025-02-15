package install

import (
	"fmt"
	"os"
	"path/filepath"
	"tachyon/internal/archive"
	"tachyon/internal/env"
)

func Package(pkgPath string, installDeps bool) error {
	fmt.Println("📦 Установка пакета:", pkgPath)

	sitePackages, err := env.GetSitePackagesPath()
	if err != nil {
		return err
	}

	err = archive.ExtractTPK(pkgPath, sitePackages)
	if err != nil {
		return err
	}

	if installDeps {
		depsPath := filepath.Join(sitePackages, filepath.Base(pkgPath)+".deps")
		if _, err := os.Stat(depsPath); err == nil {
			fmt.Println("📦 Установка зависимостей из:", depsPath)
			depsFile, err := os.ReadFile(depsPath)
			if err == nil {
				for _, dep := range filepath.SplitList(string(depsFile)) {
					fmt.Println("📦 Устанавливаем зависимость:", dep)
					Package(dep, false) // Рекурсивная установка
				}
			}
		}
	}

	return nil
}