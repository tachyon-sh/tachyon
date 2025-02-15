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

func ExtractTPK(tpkPath string, destPath string) error {
	file, err := os.Open(tpkPath)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make([]byte, 97)
	n, err := file.Read(header)
	if err != nil {
		return err
	}

	if n < 97 {
		return fmt.Errorf("❌ Ошибка: заголовок `.tpk` повреждён (ожидалось 97 байт, получено %d)", n)
	}

	expectedHash := header[:32]
	signatureLen := header[32]
	if signatureLen > 64 {
		return fmt.Errorf("❌ Ошибка: некорректная длина подписи %d (максимум 64)", signatureLen)
	}
	signature := header[33 : 33+signatureLen]

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return err
	}

	actualHash := hasher.Sum(nil)
	if !bytes.Equal(expectedHash, actualHash) {
		fmt.Printf("Ожидалось: %x\n", expectedHash)
		fmt.Printf("Получено: %x\n", actualHash)
		return fmt.Errorf("❌ Проверка целостности не пройдена! Файл повреждён")
	}

	fmt.Println("✅ SHA-256 проверен, целостность подтверждена.")

	if signatureLen == 64 {
		err := security.VerifySHA256(hex.EncodeToString(actualHash), hex.EncodeToString(signature))
		if err != nil {
			return err
		}
		fmt.Println("✅ Подпись проверена, пакет подлинный.")
	}

	_, err = file.Seek(97, io.SeekStart)
	if err != nil {
		return err
	}

	zstdReader, err := zstd.NewReader(file)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	tarReader := tar.NewReader(zstdReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		outPath := filepath.Join(destPath, header.Name)
		fmt.Println("📂 Распаковка:", outPath)

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(outPath, os.ModePerm)
			continue
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, tarReader)
		if err != nil {
			return err
		}
	}

	fmt.Println("✅ Пакет успешно установлен!")
	return nil
}