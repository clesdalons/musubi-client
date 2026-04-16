package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ZipFolder(sourceDir string) (string, error) {
	archivePath := sourceDir + ".zip"
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		entry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(entry, file)
		file.Close()
		return err
	})

	return archivePath, err
}

func ExtractZip(zipPath, targetDir, zipName string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	timestamp := time.Now().Format("20060102_150405")
	rawName := strings.TrimSuffix(zipName, filepath.Ext(zipName))
	destination := filepath.Join(targetDir, fmt.Sprintf("%s_%s", timestamp, rawName))

	for _, file := range r.File {
		filePath := filepath.Join(destination, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", filePath)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
