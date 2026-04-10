package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PullResponse correspond à la structure JSON de ton API Musubi
type PullResponse struct {
	FileName    string `json:"fileName"` // ex: "MySuperSave.zip"
	Uploader    string `json:"uploader"`
	DownloadUrl string `json:"downloadUrl"`
}

// DownloadAndExtract gère le flux complet : Metadata -> Download -> Unzip
func (a *App) DownloadAndExtract() error {
	cfg := a.LoadConfig()
	apiUrl := fmt.Sprintf("https://musubi.azurewebsites.net/api/pullsave?campaignId=%s", cfg.Campaign)

	// 1. Récupération des métadonnées
	resp, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("impossible de contacter l'API : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("erreur API : code %d", resp.StatusCode)
	}

	var meta PullResponse
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return fmt.Errorf("échec décodage JSON : %v", err)
	}

	// 2. Téléchargement du fichier binaire
	saveResp, err := http.Get(meta.DownloadUrl)
	if err != nil {
		return fmt.Errorf("échec du téléchargement Storage : %v", err)
	}
	defer saveResp.Body.Close()

	if saveResp.StatusCode != 200 {
		return fmt.Errorf("le lien de téléchargement a expiré ou est invalide")
	}

	// 3. Création d'un fichier temporaire local pour le ZIP
	// On utilise meta.FileName pour garder une trace du nom original
	tmpFilePath := filepath.Join(os.TempDir(), meta.FileName)
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("création fichier temp échouée : %v", err)
	}
	defer os.Remove(tmpFilePath) // Nettoyage après extraction

	size, err := io.Copy(tmpFile, saveResp.Body)
	tmpFile.Close() // On ferme pour pouvoir le relire avec zip.OpenReader
	if err != nil {
		return fmt.Errorf("échec écriture ZIP : %v", err)
	}

	// 4. Extraction vers le dossier Story
	return a.unzip(tmpFilePath, size, cfg.SavePath, meta.FileName)
}

// unzip extrait le contenu dans Story/NomDeLaSave/
func (a *App) unzip(zipPath string, size int64, storyPath string, zipName string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// On détermine le nom du dossier de sauvegarde (on retire le .zip)
	// "TestSave.zip" devient "TestSave"
	timestamp := time.Now().Format("20060102_150405")
	rawName := strings.TrimSuffix(zipName, filepath.Ext(zipName))
	// Résultat: 20260410_223015_AutoSave_0
	saveFolderName := fmt.Sprintf("%s_%s", timestamp, rawName)

	// Le chemin final sera par exemple : .../Documents/.../Story/TestSave/
	finalTargetDir := filepath.Join(storyPath, saveFolderName)

	for _, f := range r.File {
		// On construit le chemin du fichier à extraire
		// Si le ZIP contient déjà un dossier, f.Name sera "Dossier/Fichier.lsv"
		// Sinon ce sera juste "Fichier.lsv"
		fpath := filepath.Join(finalTargetDir, f.Name)

		// Vérification de sécurité (Zip Slip)
		if !strings.HasPrefix(fpath, filepath.Clean(finalTargetDir)+string(os.PathSeparator)) && fpath != finalTargetDir {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Création des répertoires parents
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Création du fichier destination
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Fermeture propre pour chaque fichier
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
