package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

type Favorites struct {
	IDs []string `json:"ids"`
}

func favoritesFilePath() string {
	// Ton dossier est "Data" (majuscule) chez toi.
	// Windows s’en fiche, mais on reste cohérent.
	return filepath.Join("Data", "favorites.json")
}

func LoadFavorites() Favorites {
	path := favoritesFilePath()

	// Si le fichier n’existe pas, on retourne vide.
	if _, err := os.Stat(path); err != nil {
		return Favorites{IDs: []string{}}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return Favorites{IDs: []string{}}
	}

	var favs Favorites
	if err := json.Unmarshal(b, &favs); err != nil {
		return Favorites{IDs: []string{}}
	}

	if favs.IDs == nil {
		favs.IDs = []string{}
	}
	return favs
}

func SaveFavorites(favs Favorites) error {
	path := favoritesFilePath()

	_ = os.MkdirAll(filepath.Dir(path), 0755)

	b, err := json.MarshalIndent(favs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func IsFavorite(ids []string, id string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func ToggleFavorite(id string) error {
	favs := LoadFavorites()
	if IsFavorite(favs.IDs, id) {
		// remove
		newIDs := []string{}
		for _, v := range favs.IDs {
			if v != id {
				newIDs = append(newIDs, v)
			}
		}
		favs.IDs = newIDs
	} else {
		favs.IDs = append(favs.IDs, id)
	}
	return SaveFavorites(favs)
}

func FavoriteIDsToInts(ids []string) []int {
	out := []int{}
	for _, s := range ids {
		n, err := strconv.Atoi(s)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}
