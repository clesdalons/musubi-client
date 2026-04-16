package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/musubi-client/musubi-client/internal/config"
	"github.com/musubi-client/musubi-client/internal/model"
)

type AzureClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *AzureClient {
	return &AzureClient{
		BaseURL: "https://musubi.azurewebsites.net/api",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AzureClient) GetLatestSaveInfo(campaign string) (map[string]interface{}, error) {
	if campaign == "" {
		return nil, fmt.Errorf("campaign id is required")
	}

	url := fmt.Sprintf("%s/GetLatestSaveInfo?campaignId=%s", c.BaseURL, campaign)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloud status request failed: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AzureClient) DownloadLatestSave(campaign string) (*model.PullResponse, string, error) {
	if campaign == "" {
		return nil, "", fmt.Errorf("campaign id is required")
	}

	metaURL := fmt.Sprintf("%s/pullsave?campaignId=%s", c.BaseURL, campaign)
	resp, err := c.HTTPClient.Get(metaURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("pull request failed: %s", resp.Status)
	}

	var payload model.PullResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, "", err
	}

	zipResp, err := c.HTTPClient.Get(payload.DownloadURL)
	if err != nil {
		return nil, "", err
	}
	defer zipResp.Body.Close()

	if zipResp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download failed: %s", zipResp.Status)
	}

	tempFile, err := os.CreateTemp("", "musubi-*.zip")
	if err != nil {
		return nil, "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, zipResp.Body); err != nil {
		os.Remove(tempFile.Name())
		return nil, "", err
	}

	return &payload, tempFile.Name(), nil
}

func (c *AzureClient) UploadSave(zipPath string, cfg config.Config) error {
	if cfg.Campaign == "" || cfg.Uploader == "" {
		return fmt.Errorf("campaign and uploader are required")
	}

	fileName := filepath.Base(zipPath)
	url := fmt.Sprintf(
		"%s/pushsave?campaignId=%s&uploader=%s&fileName=%s",
		c.BaseURL,
		cfg.Campaign,
		cfg.Uploader,
		fileName,
	)

	file, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodPost, url, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Save-Name", fileName)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s %s", resp.Status, string(body))
	}

	return nil
}
