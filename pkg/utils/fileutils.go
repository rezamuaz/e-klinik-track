package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SaveFileFromURL(url string, folder string) error {

	basePath, err := os.Getwd() // Get current working directory
	if err != nil {
		return err
	}

	// Navigate up one level from the base directory
	mntMangaPath := filepath.Join(basePath, "..", "mnt", "pvsave")

	// Create the mnt/manga directory if it doesn't exist
	if err := os.MkdirAll(mntMangaPath, os.ModePerm); err != nil {
		return err
	}

	// Create the custom folder inside mnt/manga
	customFolderPath := filepath.Join(mntMangaPath, folder)
	if err := os.MkdirAll(customFolderPath, os.ModePerm); err != nil {
		return err
	}

	filename := filepath.Base(url)

	myFilesPath := "/mnt/manga"
	// Create the file
	absoluteFilePath := filepath.Join(basePath, "..", myFilesPath, folder, filename)
	out, err := os.Create(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}
	time.Sleep(2 * time.Second)
	return nil
}
