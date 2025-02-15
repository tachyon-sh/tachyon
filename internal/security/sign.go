package security

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
)

func SignSHA256(hash string) (string, error) {
	privData, err := os.ReadFile("tachyon.priv")
	if err != nil {
		return "", fmt.Errorf("❌ Приватный ключ не найден, подпись пропущена")
	}

	privKey, err := hex.DecodeString(string(privData))
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(privKey, []byte(hash))
	return hex.EncodeToString(signature), nil
}

func VerifySHA256(hash string, signature string) error {
	pubData, err := os.ReadFile("tachyon.pub")
	if err != nil {
		return fmt.Errorf("🔓 Публичный ключ не найден, подпись пропущена")
	}

	pubKey, err := hex.DecodeString(string(pubData))
	if err != nil {
		return err
	}

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}

	if !ed25519.Verify(pubKey, []byte(hash), sigBytes) {
		return fmt.Errorf("❌ Подпись недействительна! Файл мог быть подменён")
	}

	return nil
}