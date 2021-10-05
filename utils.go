package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func gethn() string {
	hn := os.Getenv("MAIL_HOSTNAME")
	if len(strings.TrimSpace(hn)) > 0 {
		return hn
	}
	hn, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	return hn
}

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	if len(strings.TrimSpace(value)) > 0 {
		log.Println(name, value)
		return value
	}
	log.Println(name, defval)
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
