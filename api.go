package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const groupieAPI = "https://groupietrackers.herokuapp.com/api"

var httpClient = &http.Client{Timeout: 15 * time.Second}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`

	LocationsURL    string `json:"locations"`
	ConcertDatesURL string `json:"concertDates"`
	RelationsURL    string `json:"relations"`
}

func GetArtists() ([]Artist, error) {
	url := fmt.Sprintf("%s/artists", groupieAPI)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func GetArtistByID(id int) (Artist, error) {
	url := fmt.Sprintf("%s/artists/%d", groupieAPI, id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return Artist{}, err
	}
	defer resp.Body.Close()

	var artist Artist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return Artist{}, err
	}
	return artist, nil
}
