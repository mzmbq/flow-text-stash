package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mzmbq/flow-text-stash/internal/stash"
	"github.com/mzmbq/flow-text-stash/internal/store"
	"github.com/mzmbq/flow-text-stash/internal/utils"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ts <query>")
	}

	dataPath := filepath.Join(utils.GetDataDir(), "data.yaml")
	s, err := store.New(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	ts := stash.New(s)
	if err := ts.HandleRPC(os.Args[1]); err != nil {
		log.Fatal(err.Error())
	}
}
