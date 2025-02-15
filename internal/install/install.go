package install

import (
	"fmt"
	"tachyon/internal/archive"
	"tachyon/internal/cache"
	"tachyon/internal/env"
	"tachyon/internal/timer"
	"path/filepath" 
)

func Package(pkgPath string, installDeps bool, channel string) error {
	t := timer.StartTask(fmt.Sprintf("Установка пакета: %s (канал: %s)", pkgPath, channel))

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

	t.Stop()
	return nil
}