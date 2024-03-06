package main

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Fonction pour calculer la similarité entre deux chaînes
func similarity(query, name string) int {
	return len(name) - len(strings.ReplaceAll(strings.ToLower(name), strings.ToLower(query), ""))
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Groupie Tracker")
	window.Resize(fyne.NewSize(800, 600))

	artistNames := []string{
		"Vincent van Gogh", "Pablo Picasso", "Leonardo da Vinci", "Claude Monet",
		"Michelangelo", "Frida Kahlo", "Georgia O'Keeffe", "Salvador Dalí",
		"Andy Warhol", "Jackson Pollock", "Henri Matisse", "Wassily Kandinsky",
		"Gustav Klimt", "Edvard Munch", "René Magritte", "Marc Chagall",
		"Joan Miró", "Paul Cézanne", "Rembrandt van Rijn", "Banksy",
		"Keith Haring", "Jean-Michel Basquiat", "Yayoi Kusama", "Cindy Sherman",
		"Ai Weiwei", "Damien Hirst", "Anish Kapoor", "Jeff Koons",
		"Tracey Emin", "Marina Abramović", "Eminem",
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for an artist")

	suggestionBox := container.NewVBox()

	updateSuggestions := func(query string) {
		suggestionBox.Objects = nil

		if query != "" {
			// Trie les suggestions en fonction de leur similarité avec la requête
			sortedArtists := make([]string, len(artistNames))
			copy(sortedArtists, artistNames)
			sort.Slice(sortedArtists, func(i, j int) bool {
				return similarity(query, sortedArtists[i]) > similarity(query, sortedArtists[j])
			})

			// Ajoute les suggestions à la boîte de suggestions jusqu'au maximum possible en fonction de leur similarité
			maxSuggestions := 5
			for _, name := range sortedArtists {
				if similarity(query, name) <= 0 {
					break
				}
				suggestionBox.Add(widget.NewLabel(name))
				maxSuggestions--
				if maxSuggestions == 0 {
					break
				}
			}
		}

		window.Content().Refresh()
	}

	searchEntry.OnChanged = func(query string) {
		updateSuggestions(query)
	}

	containerBox := container.NewVBox(searchEntry, suggestionBox)

	window.SetContent(containerBox)
	window.SetIcon(theme.VolumeUpIcon())
	window.ShowAndRun()
}
