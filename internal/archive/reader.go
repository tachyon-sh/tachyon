package archive

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"tachyon/internal/security"
	"tachyon/internal/timer"

	"github.com/klauspost/compress/zstd"
	"golang.org/x/sys/unix"
)

func ExtractTPK(tpkPath string, destPath string) error {
	t := timer.StartTask(fmt.Sprintf("Распаковка .tpk: %s", tpkPath))

	file, err := os.Open(tpkPath)
	if err != nil {
		return fmt.Errorf("❌ Ошибка открытия .tpk: %w", err)
	}
	defer file.Close()

	header := make([]byte, 97)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("❌ Ошибка чтения заголовка: %w", err)
	}

	expectedHash := header[:32]
	signature := header[33:97]

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("❌ Ошибка получения информации о файле: %w", err)
	}
	mmapData, err := unix.Mmap(int(file.Fd()), 97, int(stat.Size()-97), unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return fmt.Errorf("❌ Ошибка при `mmap`: %w", err)
	}
	defer unix.Munmap(mmapData)

	fmt.Println("🔍 Проверка SHA-256...")
	hasher := sha256.New()
	hasher.Write(mmapData)
	actualHash := hasher.Sum(nil)

	if !bytes.Equal(expectedHash, actualHash) {
		return fmt.Errorf("❌ Ошибка: SHA-256 не совпадает!\nОжидалось: %x\nПолучено: %x", expectedHash, actualHash)
	}
	fmt.Println("✅ SHA-256 проверен, целостность подтверждена.")

	fmt.Println("🔍 Проверка подписи...")
	if err := security.VerifySHA256(hex.EncodeToString(actualHash), hex.EncodeToString(signature)); err != nil {
		return fmt.Errorf("❌ Ошибка подписи: %v", err)
	}
	fmt.Println("✅ Подпись подтверждена, файл подлинный.")

	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("❌ Ошибка создания директории: %w", err)
	}

	fmt.Println("📂 Извлечение файлов...")
	zstdReader, err := zstd.NewReader(bytes.NewReader(mmapData))
	if err != nil {
		return fmt.Errorf("❌ Ошибка инициализации ZSTD: %w", err)
	}
	defer zstdReader.Close()
	tarReader := tar.NewReader(zstdReader)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	fileChan := make(chan *tar.Header, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for header := range fileChan {
				outPath := filepath.Join(destPath, header.Name)
				if header.Typeflag == tar.TypeDir {
					os.MkdirAll(outPath, os.ModePerm)
					continue
				}

				outFile, err := os.Create(outPath)
				if err != nil {
					fmt.Println("❌ Ошибка создания файла:", err)
					continue
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, tarReader)
				if err != nil {
					fmt.Println("❌ Ошибка копирования:", err)
				}
			}
		}()
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("❌ Ошибка при чтении TAR: %w", err)
		}
		fileChan <- header
	}

	close(fileChan)
	wg.Wait()

	fmt.Println("✅ Пакет успешно установлен!")
	t.Stop()
	return nil
}