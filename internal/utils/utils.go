package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Wrap(s string) string {
	s = strings.TrimSpace(s)
	s = strings.SplitN(s, "\n", 2)[0]
	s = strings.TrimSpace(s)

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
