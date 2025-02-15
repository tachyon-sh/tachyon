package cache

import (
	"os"
	"path/filepath"
)

func ListCache() ([]string, error) {
	cacheDir := GetCacheDir()
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	var packages []string

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".tpk" {
			packages = append(packages, file.Name()+" (default)")
		}
	}

	channels := []string{"stable", "beta", "nightly"}
	for _, channel := range channels {
		channelDir := GetChannelDir(channel)
		channelFiles, err := os.ReadDir(channelDir)
		if err != nil {
			continue 
		}

		for _, file := range channelFiles {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".tpk" {
				packages = append(packages, file.Name()+" ("+channel+")")
			}
		}
	}

	return packages, nil
}