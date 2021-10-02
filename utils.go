package main

import (
	"os"
	"path/filepath"
	"strings"
)

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	if len(strings.TrimSpace(value)) > 0 {
		return value
	}
	return defval
}

func withext(ext string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	file := base + "." + ext
	return filepath.Join(dir, file), nil
}
