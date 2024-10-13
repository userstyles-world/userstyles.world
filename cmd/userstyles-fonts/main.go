package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	releaseURL = "https://api.github.com/repos/rsms/inter/releases/44888689" // v3.19
	distDir    = path.Join("web", "static", "fonts")
)

type (
	Release struct {
		Assets []Asset `json:"assets"`
	}

	Asset struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
		Path        string
	}
)

func ts() string {
	return time.Now().Format("[15:04:05]")
}

func getAsset() (Asset, error) {
	resp, err := http.Get(releaseURL)
	if err != nil {
		return Asset{}, err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, err
	}

	var release Release
	if err := json.Unmarshal(content, &release); err != nil {
		return Asset{}, err
	}

	if len(release.Assets) == 0 {
		return Asset{}, fmt.Errorf("failed to get an asset")
	}

	asset := release.Assets[0]
	asset.Path = path.Join("data", asset.Name)

	return asset, nil
}

func getArchive(asset Asset) error {
	// Skip re-downloading archive if it exists.
	if _, err := os.Stat(asset.Path); err == nil {
		fmt.Printf("%s Archive %q already exists.\n", ts(), asset.Name)
		return nil
	}

	fmt.Printf("%s Downloading %q archive.\n", ts(), asset.Name)
	archive, err := http.Get(asset.DownloadURL)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(archive.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(asset.Path, content, 0o755); err != nil {
		return err
	}

	return os.RemoveAll(distDir)
}

func extractArchive(asset Asset) error {
	// Skip re-extracting archive if it exists.
	if _, err := os.Stat(distDir); err == nil {
		fmt.Printf("%s Fonts are up to date.\n", ts())
		return nil
	}

	// Create dist directory if doesn't exist.
	if _, err := os.Stat(distDir); err != nil {
		if err := os.MkdirAll(distDir, 0o755); err != nil {
			return err
		}
	}

	// Open the archive file.
	archive, err := zip.OpenReader(asset.Path)
	if err != nil {
		return err
	}

	// Extract and save a file.
	extract := func(file *zip.File, output string) error {
		flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		res, err := os.OpenFile(output, flags, file.Mode())
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(res, src); err != nil {
			return err
		}

		return nil
	}

	// Extract specific fonts.
	font := func(kind string) bool {
		fonts := []string{"Regular", "Bold", "Italic", "BoldItalic"}
		for _, font := range fonts {
			if strings.HasPrefix(kind, "Inter Web/Inter-"+font) {
				return true
			}
		}

		return false
	}

	// Iterate over files and extract the ones we want.
	for _, file := range archive.File {
		if font(file.Name) {
			name := path.Base(file.Name)
			output := path.Join(distDir, name)
			if strings.HasPrefix(output, filepath.Clean(output)) {
				fmt.Printf("%s Extracting %q to %q.\n", ts(), name, output)
				if err := extract(file, output); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func main() {
	asset, err := getAsset()
	if err != nil {
		log.Fatal(err)
	}

	if err := getArchive(asset); err != nil {
		log.Fatal(err)
	}

	if err := extractArchive(asset); err != nil {
		log.Fatal(err)
	}
}
