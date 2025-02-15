package install

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemovePackage(pkgName string) error {
	sitePackages, err := filepath.Abs("/usr/local/lib/python3.10/site-packages") // Заменить на динамическое получение пути
	if err != nil {
		return err
	}

	pkgPath := filepath.Join(sitePackages, pkgName)
	infoPath := filepath.Join(sitePackages, pkgName+".dist-info")

	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return fmt.Errorf("пакет %s не найден", pkgName)
	}

	err = os.RemoveAll(pkgPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(infoPath)
	if err != nil {
		return err
	}

	return nil
}