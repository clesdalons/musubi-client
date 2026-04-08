package main

import (
	"context"
	"os"
	"path/filepath"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetSavePath retourne le chemin vers les sauvegardes
func (a *App) GetSavePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "Erreur : Home introuvable"
	}

	path := filepath.Join(home, "Documents", "Larian Studios", "Divinity Original Sin 2 Enhanced Edition", "Savegames", "Story")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "Dossier DOS2 non trouvé"
	}
	return path
}
