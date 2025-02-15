package install

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func GetSitePackagesPath() (string, error) {
	cmd := exec.Command("python3", "-c", "import site; print(site.getsitepackages()[0])")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("не удалось определить site-packages через Python: %w", err)
	}

	sitePackages := strings.TrimSpace(out.String())
	absPath, err := filepath.Abs(sitePackages)
	if err != nil {
		return "", fmt.Errorf("ошибка получения абсолютного пути site-packages: %w", err)
	}

	return absPath, nil
}

func RemovePackage(pkgName string) error {
	sitePackages, err := GetSitePackagesPath()
	if err != nil {
		return err
	}

	pkgPath := filepath.Join(sitePackages, pkgName)
	infoPath := filepath.Join(sitePackages, pkgName+".dist-info")

	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return fmt.Errorf("пакет %s не найден", pkgName)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := os.RemoveAll(pkgPath); err != nil {
			fmt.Printf("⚠️ Ошибка удаления %s: %v\n", pkgPath, err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := os.RemoveAll(infoPath); err != nil {
			fmt.Printf("⚠️ Ошибка удаления %s.dist-info: %v\n", infoPath, err)
		}
	}()

	wg.Wait()
	fmt.Printf("✅ Пакет %s успешно удалён!\n", pkgName)
	return nil
}

func RemovePackageFromCache(pkgName, channel string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
	}

	cachePath := filepath.Join(homeDir, ".tachyon", "cache", "channels", channel, pkgName+".tpk")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil 
	}

	if err := os.Remove(cachePath); err != nil {
		return fmt.Errorf("ошибка удаления кеша %s: %w", cachePath, err)
	}

	fmt.Printf("🗑️  Пакет %s удалён из кеша (%s)\n", pkgName, channel)
	return nil
}

func RemovePackageAll(pkgName string) error {
	if err := RemovePackage(pkgName); err != nil {
		return err
	}

	channels := []string{"stable", "beta", "nightly"}
	var wg sync.WaitGroup

	for _, channel := range channels {
		wg.Add(1)
		go func(ch string) {
			defer wg.Done()
			_ = RemovePackageFromCache(pkgName, ch)
		}(channel)
	}

	wg.Wait()
	return nil
}