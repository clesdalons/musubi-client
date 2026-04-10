// uploader.go
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ZipFolder creates a zip archive of the specified save directory
func (a *App) ZipFolder(sourceDir string) (string, error) {
	archivePath := sourceDir + ".zip"
	archive, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	defer writer.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create the file header in the zip
		f, err := writer.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		return err
	})

	return archivePath, err
}

func (a *App) UploadToAzure(zipPath string) error {
	cfg := a.GetSettings()
	fileName := filepath.Base(zipPath)

	// Build the dynamic URL with parameters
	url := fmt.Sprintf(
		"https://musubi.azurewebsites.net/api/pushsave?campaignId=%s&uploader=%s&fileName=%s",
		cfg.Campaign,
		cfg.Uploader,
		fileName,
	)

	log.Printf("[Uploader] Target URL: %s", url)

	file, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(zipPath)

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Save-Name", filepath.Base(zipPath))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	return nil
}
