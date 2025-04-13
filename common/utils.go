package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

func Download(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http request GET '%s' failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := WriteFile(path, resp.Body, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

func ReadFile(filename string, w io.Writer) error {
	path, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to convert '%s' to a absolute path: %w", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", filename, err)
	}

	return nil
}

func WriteFile(filename string, r io.Reader, perm os.FileMode) error {
	path, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to convert '%s' to a absolute path: %w", filename, err)
	}

	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0644); err != nil {
		return fmt.Errorf("failed to create parent directory '%s': %w", parent, err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create '%s': %w", filename, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", filename, err)
	}

	return nil
}

func DownloadArtifact(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http request GET '%s' failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := WriteFile("extensions/marketplace.json", resp.Body, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func LoadJsonFS(root fs.FS, path string, value any) error {
	f, err := root.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(value); err != nil {
		return fmt.Errorf("failed to json decode '%s': %w", path, err)
	}

	return nil
}

func DownloadJson(url string, value any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http request GET '%s' failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(value); err != nil {
		return fmt.Errorf("failed to json decode: %w", err)
	}

	return nil
}
