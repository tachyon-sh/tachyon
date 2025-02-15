package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func PackTPK(sourceDir, outputFile string) error {
	fmt.Println("📦 Упаковка пакета:", sourceDir)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("❌ Ошибка создания .tpk: %w", err)
	}
	defer outFile.Close()

	zstdWriter, err := zstd.NewWriter(outFile)
	if err != nil {
		return fmt.Errorf("❌ Ошибка инициализации ZSTD: %w", err)
	}

	tarWriter := tar.NewWriter(zstdWriter)

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		finalPath := filepath.Join("test-package", relPath)

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = finalPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("❌ Ошибка записи заголовка TAR: %w", err)
		}

		src, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("❌ Ошибка открытия файла: %w", err)
		}
		defer src.Close()

		if _, err := io.Copy(tarWriter, src); err != nil {
			return fmt.Errorf("❌ Ошибка копирования в TAR: %w", err)
		}

		fmt.Println("📦 Добавлен в архив:", finalPath)
		return nil
	})
	if err != nil {
		return fmt.Errorf("❌ Ошибка обхода файлов: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("❌ Ошибка закрытия TAR: %w", err)
	}
	if err := zstdWriter.Close(); err != nil {
		return fmt.Errorf("❌ Ошибка закрытия ZSTD: %w", err)
	}

	stat, err := outFile.Stat()
	if err != nil {
		return fmt.Errorf("❌ Ошибка Stat() .tpk: %w", err)
	}
	fmt.Printf("📦 Упаковка завершена, размер архива: %d байт\n", stat.Size())

	return nil
}