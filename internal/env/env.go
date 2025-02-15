package env

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPythonVersion() (string, error) {
	cmd := exec.Command("python3", "-c", "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", errors.New("не удалось определить версию Python")
	}
	return strings.TrimSpace(out.String()), nil
}

func GetSitePackagesPath() (string, error) {
	pythonVersion, err := GetPythonVersion()
	if err != nil {
		return "", err
	}

	if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
		return filepath.Join(venv, "lib", "python"+pythonVersion, "site-packages"), nil
	}

	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return filepath.Join("/opt/homebrew/lib/python"+pythonVersion, "site-packages"), nil
		} else {
			return filepath.Join("/usr/local/lib/python"+pythonVersion, "site-packages"), nil
		}
	}

	return filepath.Join("/usr/local/lib/python"+pythonVersion, "site-packages"), nil
}