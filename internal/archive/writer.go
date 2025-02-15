package archive

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"bytes"
	"tachyon/internal/security"

	"github.com/klauspost/compress/zstd"
)

func CreateTPK(inputDir string, outputFile string) error {
	tempFile := outputFile + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer out.Close()

	headerPlaceholder := make([]byte, 97) // 32 (SHA-256) + 1 (длина подписи) + 64 (подпись)
	_, err = out.Write(headerPlaceholder)
	if err != nil {
		return err
	}

	var archiveBuffer bytes.Buffer
	hasher := sha256.New()
	multiWriter := io.MultiWriter(&archiveBuffer, hasher)

	zstdWriter, err := zstd.NewWriter(multiWriter)
	if err != nil {
		return err
	}
	defer zstdWriter.Close()

	tarWriter := tar.NewWriter(zstdWriter)
	defer tarWriter.Close()

	err = filepath.Walk(inputDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(inputDir, filePath)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		fmt.Println("Добавлен в архив:", relPath)
		return nil
	})

	if err != nil {
		return err
	}

	tarWriter.Close()
	zstdWriter.Close()

	hash := hasher.Sum(nil)

	signature, err := security.SignSHA256(hex.EncodeToString(hash))
	signatureBytes := make([]byte, 64) // Гарантируем, что подпись всегда 64 байта
	signatureLen := byte(0)

	if err == nil {
		sigDecoded, _ := hex.DecodeString(signature)
		copy(signatureBytes[:], sigDecoded) // Заполняем `signatureBytes` правильными данными
		signatureLen = 64
		fmt.Println("✅ Файл подписан.")
	} else {
		fmt.Println("⚠️ Подпись отсутствует, продолжаем без неё.")
	}

	_, err = out.WriteAt(hash, 0)                           // 32 байта SHA-256
	_, err = out.WriteAt([]byte{signatureLen}, 32)         // 1 байт длины подписи
	_, err = out.WriteAt(signatureBytes[:signatureLen], 33) // 64 байта подписи (или пустой массив)

	_, err = io.Copy(out, &archiveBuffer)
	if err != nil {
		return err
	}

	err = os.Rename(tempFile, outputFile)
	if err != nil {
		return err
	}

	fmt.Println("📦 Пакет успешно упакован:", outputFile)
	return nil
}