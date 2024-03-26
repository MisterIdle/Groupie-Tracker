package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type apiCache struct {
	artistsCache   map[string][]Artist
	locationsCache map[int][]string
}

var cache = &apiCache{
	artistsCache:   make(map[string][]Artist),
	locationsCache: make(map[int][]string),
}

func fetchArtists() ([]Artist, error) {
	if artists, ok := cache.artistsCache["artists"]; ok {
		return artists, nil
	}

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	cache.artistsCache["artists"] = artists

	return artists, nil
}

func fetchLocations(artistID int) ([]string, error) {
	if locations, ok := cache.locationsCache[artistID]; ok {
		return locations, nil
	}

	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", artistID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locationData struct {
		Locations []string `json:"locations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&locationData); err != nil {
		return nil, err
	}

	cache.locationsCache[artistID] = locationData.Locations

	return locationData.Locations, nil
}

func preloadData() {
	_, err := fetchArtists()
	if err != nil {
		fmt.Println("Erreur lors du préchargement des artistes:", err)
	}

	artists, err := fetchArtists()
	if err != nil {
		fmt.Println("Erreur lors du préchargement des emplacements:", err)
	}
	for _, artist := range artists {
		_, err := fetchLocations(artist.ID)
		if err != nil {
			fmt.Println("Erreur lors du préchargement des emplacements pour l'artiste", artist.ID, ":", err)
		}
	}
}

func init() {
	preloadData()
}
