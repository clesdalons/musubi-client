package model

type PullResponse struct {
	FileName    string `json:"fileName"`
	Uploader    string `json:"uploader"`
	DownloadURL string `json:"downloadUrl"`
}
