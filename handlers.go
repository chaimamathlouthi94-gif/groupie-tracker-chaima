package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type FilmsPageData struct {
	Title      string
	Artists    []Artist
	Favs       []string
	Page       int
	Limit      int
	TotalPages int
	PrevPage   int
	NextPage   int
}

type FilmDetailData struct {
	Title  string
	Artist Artist
	Favs   []string
	IsFav  bool
}

type FavoritesPageData struct {
	Title   string
	Artists []Artist
	Favs    []string
}

func clampLimit(n int) int {
	// Sujet: 10, 20, 30
	if n == 20 || n == 30 {
		return n
	}
	return 10
}

func FilmsHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artists, err := GetArtists()
		if err != nil {
			http.Error(w, "Erreur API artists", http.StatusInternalServerError)
			return
		}

		favs := LoadFavorites()

		// pagination params
		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				page = n
			}
		}
		limit := 10
		if l := r.URL.Query().Get("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil {
				limit = clampLimit(n)
			}
		} else {
			limit = 10
		}

		total := len(artists)
		totalPages := total / limit
		if total%limit != 0 {
			totalPages++
		}
		if totalPages == 0 {
			totalPages = 1
		}
		if page > totalPages {
			page = totalPages
		}

		start := (page - 1) * limit
		end := start + limit
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		pageArtists := artists[start:end]

		prev := page - 1
		if prev < 1 {
			prev = 1
		}
		next := page + 1
		if next > totalPages {
			next = totalPages
		}

		data := FilmsPageData{
			Title:      "Artists",
			Artists:    pageArtists,
			Favs:       favs.IDs,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			PrevPage:   prev,
			NextPage:   next,
		}

		if err := tpl.ExecuteTemplate(w, "films.html", data); err != nil {
			http.Error(w, "Erreur template films", http.StatusInternalServerError)
			return
		}
	}
}

func FilmDetailHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URL: /films/12
		idStr := strings.TrimPrefix(r.URL.Path, "/films/")
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			http.Redirect(w, r, "/films", http.StatusFound)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		artist, err := GetArtistByID(id)
		if err != nil {
			http.Error(w, "Erreur API artist detail", http.StatusInternalServerError)
			return
		}

		favs := LoadFavorites()
		idAsString := strconv.Itoa(id)

		data := FilmDetailData{
			Title:  artist.Name,
			Artist: artist,
			Favs:   favs.IDs,
			IsFav:  IsFavorite(favs.IDs, idAsString),
		}

		if err := tpl.ExecuteTemplate(w, "film.html", data); err != nil {
			http.Error(w, "Erreur template film", http.StatusInternalServerError)
			return
		}
	}
}

func FavoritesHandler(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		favs := LoadFavorites()
		all, err := GetArtists()
		if err != nil {
			http.Error(w, "Erreur API artists", http.StatusInternalServerError)
			return
		}

		// filtre: on garde ceux qui sont favoris
		favSet := map[string]bool{}
		for _, id := range favs.IDs {
			favSet[id] = true
		}

		list := []Artist{}
		for _, a := range all {
			if favSet[strconv.Itoa(a.ID)] {
				list = append(list, a)
			}
		}

		data := FavoritesPageData{
			Title:   "Favoris",
			Artists: list,
			Favs:    favs.IDs,
		}

		if err := tpl.ExecuteTemplate(w, "favorites.html", data); err != nil {
			http.Error(w, "Erreur template favorites", http.StatusInternalServerError)
			return
		}
	}
}

func ToggleFavoriteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/films", http.StatusFound)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			http.Redirect(w, r, "/films", http.StatusFound)
			return
		}

		_ = ToggleFavorite(id)

		// Retour à la page précédente si possible
		ref := r.Referer()
		if ref == "" {
			ref = "/films"
		}
		http.Redirect(w, r, ref, http.StatusFound)
	}
}
