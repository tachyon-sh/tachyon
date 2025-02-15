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
	"runtime"
	"sync"
	"tachyon/internal/security"

	"github.com/klauspost/compress/zstd"
	"golang.org/x/sys/unix"
)

func ExtractTPK(tpkPath string, destPath string) error {
	file, err := os.Open(tpkPath)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make([]byte, 97)
	_, err = file.Read(header)
	if err != nil {
		return err
	}

	expectedHash := header[:32]
	signatureLen := header[32]
	signature := header[33 : 33+signatureLen]

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return err
	}

	actualHash := hasher.Sum(nil)
	if !bytes.Equal(expectedHash, actualHash) {
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

	file.Seek(97, io.SeekStart)
	fd := int(file.Fd())
	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	mmapData, err := unix.Mmap(fd, 97, int(fileSize-97), unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return err
	}

	zstdReader, err := zstd.NewReader(bytes.NewReader(mmapData))
	if err != nil {
		return err
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
					fmt.Println("Ошибка:", err)
					continue
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, tarReader)
				if err != nil {
					fmt.Println("Ошибка копирования:", err)
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
			return err
		}
		fileChan <- header
	}

	close(fileChan)
	wg.Wait()

	fmt.Println("✅ Пакет успешно установлен!")
	return nil
}