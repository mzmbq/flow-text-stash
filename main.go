package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/atotto/clipboard"
	"github.com/goccy/go-yaml"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mzmbq/flow-launcher-go"
)

type Store struct {
	Path  string
	Store map[string]string
	Keys  []string
}

var store *Store

func NewStore(path string) (*Store, error) {
	s := &Store{
		Path:  path,
		Store: make(map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, &s.Store)
	if err != nil {
		return nil, err
	}

	for k := range s.Store {
		s.Keys = append(s.Keys, k)
	}

	return s, nil
}

// Look up key using fuzzy search
func (s *Store) GetFuzzy(key string) []string {
	matches := fuzzy.RankFindFold(key, s.Keys)

	sort.Sort(matches)

	var result []string
	for _, r := range matches {
		result = append(result, r.Target)
	}
	return result
}

func (s *Store) Set(key, value string) {
	s.Store[key] = value
	if !slices.Contains(s.Keys, key) {
		s.Keys = append(s.Keys, key)
	}
}

func (s *Store) Save() error {
	data, err := yaml.Marshal(s.Store)
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.Path, data, 0644); err != nil {
		return err
	}
	return nil
}

func wrap(s string) string {
	if len(s) < 20 {
		return s
	}
	return fmt.Sprintf("%s...", s[0:17])
}

func listAllStashes(req *flow.Request) *flow.Response {
	res := flow.NewResponse(req)
	for k, v := range store.Store {
		res.AddResult(&flow.Result{
			Title:    k,
			SubTitle: wrap(v),
			IcoPath:  "paste.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "paste",
				Parameters: []string{v},
			},
		})
	}
	return res
}

func handleQuery(req *flow.Request) *flow.Response {
	res := flow.NewResponse(req)
	target := req.Parameters[0]
	if target == "" {
		return listAllStashes(req)
	}
	matches := store.GetFuzzy(target)

	// List matches
	for _, m := range matches {
		val := store.Store[m]
		res.AddResult(&flow.Result{
			Title:    m,
			SubTitle: wrap(val),
			IcoPath:  "paste.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "paste",
				Parameters: []string{val},
			},
		})
	}

	// List create option for when no matches or no exact match was found
	if len(matches) == 0 || matches[0] != target {
		res.AddResult(&flow.Result{
			Title:    fmt.Sprintf("Create a paste: %s", target),
			SubTitle: "",
			IcoPath:  "add.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "create",
				Parameters: []string{target},
			},
		})
	}

	return res
}

// Create a new stash
func handleCreate(params []string) *flow.Response {
	text, err := clipboard.ReadAll()
	if err != nil {
		// TODO: handle
		return nil
	}

	store.Set(params[0], text)
	err = store.Save()
	if err != nil {
		// TODO: handle
		return nil
	}
	return nil
}

func handlePaste(params []string) *flow.Response {
	err := clipboard.WriteAll(params[0])
	if err != nil {
		// TODO: handle
		return nil
	}
	return nil
}

func getDataDir() string {
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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ts <query>")
	}

	s, err := NewStore(filepath.Join(getDataDir(), "data.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	store = s

	p := flow.NewPlugin()
	p.Query(handleQuery)
	p.Action("paste", handlePaste)
	p.Action("create", handleCreate)
	if err := p.HandleRPC(os.Args[1]); err != nil {
		log.Fatal(err.Error())
	}
}
