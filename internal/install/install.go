package install

import (
	"fmt"
	"path/filepath"
	"tachyon/internal/archive"
	"tachyon/internal/cache"
	"tachyon/internal/env"
)

func Package(pkgPath string, installDeps bool, channel string) error {
	fmt.Println("📦 Установка пакета:", pkgPath, "(канал:", channel, ")")

	sitePackages, err := env.GetSitePackagesPath()
	if err != nil {
		return err
	}

	cachedPkg := cache.GetFromChannel(filepath.Base(pkgPath), channel)
	if cache.CacheExists(filepath.Base(cachedPkg)) {
		fmt.Println("📡 Используем версию из канала:", channel, "➡", cachedPkg)
		pkgPath = cachedPkg
	} else {
		err := cache.SaveToChannel(pkgPath, channel)
		if err != nil {
			fmt.Println("⚠️ Ошибка кеширования:", err)
		}
	}

	err = archive.ExtractTPK(pkgPath, sitePackages)
	if err != nil {
		return err
	}

	return nil
}