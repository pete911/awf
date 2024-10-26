package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pete911/awf/internal/types"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	rootDir     = ".awf"
	accountFile = "_account"
)

type File struct {
	logger *slog.Logger
	dir    string
}

func LoadFile(logger *slog.Logger) (File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return File{}, err
	}

	return File{
		logger: logger.With("component", "file-store"),
		dir:    filepath.Join(home, rootDir),
	}, nil
}

func initFile(logger *slog.Logger, account types.Account, region string) (File, error) {
	f, err := LoadFile(logger)
	if err != nil {
		return File{}, err
	}

	if err := os.MkdirAll(filepath.Join(f.dir, account.Id, region), 0755); err != nil {
		return File{}, err
	}

	if err := f.write(account, "", accountFile, account); err != nil {
		return File{}, fmt.Errorf("write account data: %w", err)
	}
	return f, nil
}

func (f File) ListAccounts() ([]types.Account, error) {
	entry, err := os.ReadDir(f.dir)
	if err != nil {
		return nil, err
	}

	var accounts []types.Account
	for _, e := range entry {
		if e.IsDir() {
			var account types.Account
			if err := f.read(filepath.Join(f.dir, e.Name(), accountFile), &account); err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}

func (f File) ListRegions(account types.Account) ([]string, error) {
	entry, err := os.ReadDir(filepath.Join(f.dir, account.Id))
	if err != nil {
		return nil, err
	}

	var global bool
	var regions []string
	for _, e := range entry {
		if e.IsDir() {
			regions = append(regions, e.Name())
		}
		// we found some file that is not '_account', this means that we have global region files e.g. iam, route53 ...
		if e.Type().IsRegular() && e.Name() != accountFile {
			global = true
		}
	}
	if global {
		regions = append(regions, "")
	}
	return regions, nil
}

// Read reads and (json) marshals content of the file . Region can be empty (e.g. route53).
// If the file is not found, NotFoundError is returned.
func (f File) Read(account types.Account, region, name string, v any) error {
	path := f.filePath(account, region, name)
	return f.read(path, v)
}

func (f File) read(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewNotFoundError(fmt.Sprintf("read: %s file does not exist, empty content", path))
		}
		return err
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return nil
}

// Write writes content of the supplied (json) struct under supplied <name> file. Region
// can be empty (e.g. route53)
func (f File) write(account types.Account, region, name string, v any) error {
	path := f.filePath(account, region, name)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	return os.WriteFile(path, b, 0644)
}

func (f File) filePath(account types.Account, region, name string) string {
	if region == "" {
		return filepath.Join(f.dir, account.Id, name)
	}
	return filepath.Join(f.dir, account.Id, region, name)
}
