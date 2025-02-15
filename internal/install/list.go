package install

import (
	"fmt"
	"os"
	"strings"
)

func ListPackages(sitePackages string) ([]string, error) {
	files, err := os.ReadDir(sitePackages)
	if err != nil {
		return nil, err
	}

	var packages []string
	for _, file := range files {
		if file.IsDir() && strings.HasSuffix(file.Name(), ".dist-info") {
			pkgName := strings.TrimSuffix(file.Name(), ".dist-info")
			packages = append(packages, pkgName)
		}
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("не найдено установленных пакетов")
	}

	return packages, nil
}