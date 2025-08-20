package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafeOpen valida que filename no escape del directorio de trabajo
func SafeOpen(filename string) (*os.File, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	absPath = filepath.Clean(absPath)
	if !strings.HasPrefix(absPath, cwd+string(os.PathSeparator)) {
		return nil, fmt.Errorf("invalid imaging path: %s", filename)
	}
	return os.Open(absPath)
}
