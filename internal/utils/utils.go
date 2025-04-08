package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Wrap(s string) string {
	if len(s) < 20 {
		return s
	}
	return fmt.Sprintf("%s...", s[0:17])
}

func GetDataDir() string {
	config, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	data := filepath.Join(config, "mzmbq", "TextStash")
	if err := os.MkdirAll(data, 0755); err != nil {
		log.Fatal(err)
	}
	return data
}
