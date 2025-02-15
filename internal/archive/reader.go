package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"golang.org/x/sys/unix"
)

func ExtractTPK(tpkPath, destPath string) error {
	info, err := os.Stat(tpkPath)
	if err != nil {
		return fmt.Errorf("не могу Stat('%s'): %w", tpkPath, err)
	}
	if info.Size() < 1024 {
		// fallback
		fmt.Printf("ℹ️ Файл %s слишком маленький (%d байт), читаем напрямую без mmap\n", tpkPath, info.Size())
		return extractWithoutMmap(tpkPath, destPath)
	}
	return extractWithMmap(tpkPath, destPath)
}

func extractWithoutMmap(tpkPath, destPath string) error {
	data, err := os.ReadFile(tpkPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения .tpk: %w", err)
	}
	zr, err := zstd.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("ошибка инициализации ZSTD: %w", err)
	}
	defer zr.Close()

	tr := tar.NewReader(zr)
	return extractTar(tr, destPath)
}

func extractWithMmap(tpkPath, destPath string) error {
	f, err := os.Open(tpkPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия '%s': %w", tpkPath, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := int(stat.Size())

	mmapData, err := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return fmt.Errorf("ошибка mmap: %w", err)
	}
	defer unix.Munmap(mmapData)

	zr, err := zstd.NewReader(bytes.NewReader(mmapData))
	if err != nil {
		return fmt.Errorf("ошибка инициализации ZSTD: %w", err)
	}
	defer zr.Close()

	tr := tar.NewReader(zr)
	return extractTar(tr, destPath)
}

func extractTar(tr *tar.Reader, destPath string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка чтения TAR: %w", err)
		}

		outPath := filepath.Join(destPath, header.Name)
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return fmt.Errorf("ошибка создания директории '%s': %w", outPath, err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				return fmt.Errorf("ошибка создания директорий '%s': %w", outPath, err)
			}
			outFile, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("ошибка создания файла '%s': %w", outPath, err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("ошибка записи '%s': %w", outPath, err)
			}
			outFile.Close()
		}
	}
	fmt.Printf("✅ Распаковка успешно завершена в %s\n", destPath)
	return nil
}