package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func GenerateKeys() error {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	pubHex := hex.EncodeToString(pubKey)
	privHex := hex.EncodeToString(privKey)

	err = os.WriteFile("tachyon.pub", []byte(pubHex), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile("tachyon.priv", []byte(privHex), 0600)
	if err != nil {
		return err
	}

	fmt.Println("🔑 Ключи сгенерированы:")
	fmt.Println("📄 Публичный ключ (tachyon.pub):", pubHex)
	fmt.Println("🔐 Приватный ключ (tachyon.priv):", privHex)
	return nil
}