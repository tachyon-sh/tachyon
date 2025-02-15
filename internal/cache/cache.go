package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/tachyon-cache" 
	}
	return filepath.Join(homeDir, ".tachyon", "cache")
}

func CacheExists(pkgName string) bool {
	cacheFile := filepath.Join(GetCacheDir(), pkgName)
	_, err := os.Stat(cacheFile)
	return err == nil
}

func GetCachedFile(pkgName string) string {
	return filepath.Join(GetCacheDir(), pkgName)
}

func SaveToCache(pkgPath string) error {
	cacheDir := GetCacheDir()
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}

	destPath := filepath.Join(cacheDir, filepath.Base(pkgPath))
	input, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, input, 0644)
	if err != nil {
		return err
	}

	fmt.Println("🗂 Пакет кеширован:", destPath)
	return nil
}