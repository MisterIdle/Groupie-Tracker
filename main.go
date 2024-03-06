package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Artist représente un artiste
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

// fetchArtists récupère les données des artistes depuis l'API
func fetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// Fonction pour filtrer les artistes et les groupes
func filterArtistsAndGroups(query string, artists []Artist) []Artist {
	var filtered []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			filtered = append(filtered, artist)
		} else {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
					filtered = append(filtered, artist)
					break
				}
			}
		}
	}
	return filtered
}

// Fonction pour mettre à jour les suggestions d'artistes et de groupes
func updateSuggestions(query string, artists []Artist, suggestionBox *fyne.Container) {
	suggestionBox.Objects = nil

	if query != "" {
		filtered := filterArtistsAndGroups(query, artists)
		for _, item := range filtered {
			if len(suggestionBox.Objects) >= 6 {
				break
			}
			labelText := fmt.Sprintf("%s", item.Name)
			if len(item.Members) > 0 {
				labelText += " (" + strings.Join(item.Members, ", ") + ")"
			}
			button := widget.NewButton(labelText, func() {})
			button.Importance = widget.HighImportance
			button.Alignment = widget.ButtonAlignLeading
			button.SetIcon(theme.VolumeDownIcon())

			suggestionBox.Add(button)
		}
	}
}

// RunGroupieTracker lance l'application Groupie Tracker
func RunGroupieTracker() error {
	myApp := app.New()
	window := myApp.NewWindow("Groupie Tracker")
	window.Resize(fyne.NewSize(800, 600))

	var artists []Artist
	var err error
	if artists, err = fetchArtists(); err != nil {
		return fmt.Errorf("failed to fetch artists: %w", err)
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for an artist or a band")
	searchEntry.Resize(fyne.NewSize(400, 40)) // Ajuster la taille de la barre de recherche

	suggestionBox := container.NewVBox()
	suggestionBox.Resize(fyne.NewSize(400, 400)) // Ajuster la taille de la boîte de suggestions

	searchEntry.OnChanged = func(query string) {
		updateSuggestions(query, artists, suggestionBox)
		window.Content().Refresh()
	}

	containerBox := container.NewVBox(
		searchEntry,
		suggestionBox,
	)

	window.SetContent(containerBox)
	window.SetIcon(theme.VolumeUpIcon())
	window.ShowAndRun()

	return nil
}

func main() {
	if err := RunGroupieTracker(); err != nil {
		fmt.Println("Failed to run Groupie Tracker:", err)
	}
}
