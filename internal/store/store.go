package store

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"sort"

	"github.com/goccy/go-yaml"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Store struct {
	Path string
	Data map[string]string
	Keys []string
}

func New(path string) (*Store, error) {
	s := &Store{
		Path: path,
		Data: make(map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(data, &s.Data)
	if err != nil {
		return nil, err
	}

	for k := range s.Data {
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
	s.Data[key] = value
	if !slices.Contains(s.Keys, key) {
		s.Keys = append(s.Keys, key)
	}
}

func (s *Store) Save() error {
	data, err := yaml.Marshal(s.Data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.Path, data, 0644); err != nil {
		return err
	}
	return nil
}
