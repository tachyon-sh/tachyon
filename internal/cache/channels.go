package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetChannelDir(channel string) string {
	return filepath.Join(GetCacheDir(), "channels", channel)
}

func SaveToChannel(pkgPath string, channel string) error {
	channelDir := GetChannelDir(channel)
	err := os.MkdirAll(channelDir, os.ModePerm)
	if err != nil {
		return err
	}

	destPath := filepath.Join(channelDir, filepath.Base(pkgPath))
	input, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, input, 0644)
	if err != nil {
		return err
	}

	fmt.Println("📡 Пакет сохранён в канал:", channel, "➡", destPath)
	return nil
}

func GetFromChannel(pkgName string, channel string) string {
	return filepath.Join(GetChannelDir(channel), pkgName)
}