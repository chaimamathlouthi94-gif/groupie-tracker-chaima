package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	// 🔧 Fonctions accessibles dans les templates
	funcMap := template.FuncMap{
		"IsFavorite": IsFavorite,
	}

	// 🔧 Chargement des templates AVEC FuncMap
	tpl := template.Must(
		template.New("").Funcs(funcMap).ParseGlob("Templates/*.html"),
	)

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("Static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/films", http.StatusFound)
	})

	mux.HandleFunc("/films", FilmsHandler(tpl))
	mux.HandleFunc("/films/", FilmDetailHandler(tpl))

	mux.HandleFunc("/favorites", FavoritesHandler(tpl))
	mux.HandleFunc("/favorites/toggle", ToggleFavoriteHandler())

	addr := ":8080"
	log.Println("Serveur lancé sur http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
