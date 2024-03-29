package app

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ga *GroupieApp) searchArtists(query string) {
	if query == "all" {
		query = ""
	}

	if query == "" {
		ga.cityDropdown.Selected = "All"
	}

	filteredCards := ga.filterCards(query)

	if len(filteredCards) == 0 {
		noResultsLabel := widget.NewLabel("Aucun résultat trouvé pour la recherche : " + query)
		ga.content.Objects[1] = noResultsLabel
		ga.content.Refresh()
		ga.search.SetText("")
		return
	}
	filteredContent := container.NewVScroll(container.NewGridWithColumns(3, filteredCards...))
	ga.content.Objects[1] = filteredContent
	ga.content.Refresh()
	ga.search.SetText("")
}

func (ga *GroupieApp) updateSuggestions(query string) {
	ga.suggestionsBox.Objects = nil
	var filtered []Artist

	if query != "" {
		queryInt, err := strconv.Atoi(query)

		for _, artist := range ga.artists {
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

			if err == nil && artist.CreationDate == queryInt {
				filtered = append(filtered, artist)
			}

			if strings.Contains(artist.FirstAlbum, query) {
				filtered = append(filtered, artist)
			}
		}

		citySelected := ga.city != "all"

		allUnchecked := true
		for _, checked := range ga.checkedMembers {
			if checked {
				allUnchecked = false
				break
			}
		}

		for _, item := range filtered {
			if len(ga.suggestionsBox.Objects) >= 5 {
				break
			}

			if allUnchecked || ga.checkedMembers[len(item.Members)] || !citySelected {
				loc, err := fetchLocations(item.ID)
				if err != nil {
					fmt.Println("Erreur lors de la récupération des emplacements pour l'artiste", item.Name, ":", err)
					continue
				}

				if !citySelected || (len(loc) > 0 && containsLocation(loc, ga.city)) {
					label := item.Name
					if len(item.Members) > 0 {
						label += " (" + strings.Join(item.Members, ", ") + ")"
					}
					button := widget.NewButton(label, func(artist Artist) func() {
						return func() {
							ga.searchArtists(artist.Name)
						}
					}(item))
					button.Importance = widget.HighImportance
					button.Alignment = widget.ButtonAlignLeading
					ga.suggestionsBox.Add(button)
				}
			}
		}
	}
}

func containsLocation(locations []string, city string) bool {
	for _, loc := range locations {
		if strings.Contains(strings.ToLower(loc), strings.ToLower(city)) {
			return true
		}
	}
	return false
}

func (ga *GroupieApp) filterArtistByLocation(location string) []Artist {
	var filtered []Artist

	for _, artist := range ga.artists {
		locations, err := fetchLocations(artist.ID)
		if err != nil {
			fmt.Println("Erreur lors de la récupération des emplacements pour l'artiste", artist.Name, ":", err)
			continue
		}
		for _, loc := range locations {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(location)) {
				filtered = append(filtered, artist)
				break
			}
		}
	}

	return filtered
}

func (ga *GroupieApp) filterCards(query string) []fyne.CanvasObject {
	var filtered []fyne.CanvasObject
	queryLower := strings.ToLower(query)
	queryInt, err := strconv.Atoi(query)
	dateSearch := err == nil

	for _, artist := range ga.artists {
		allUnchecked := true
		for _, checked := range ga.checkedMembers {
			if checked {
				allUnchecked = false
				break
			}
		}

		includeArtist := strings.Contains(strings.ToLower(artist.Name), queryLower)

		// Filtrer par date si une recherche de date est spécifiée
		if dateSearch && artist.CreationDate == queryInt {
			includeArtist = true
		}

		// Filtrer par emplacement si une recherche d'emplacement est spécifiée
		if query != "" {
			locations := ga.filterArtistByLocation(query)
			for _, locArtist := range locations {
				if locArtist.ID == artist.ID {
					includeArtist = true
					break
				}
			}
		}

		// Nouvelle condition pour filtrer par date de création
		if !includeArtist && ga.creationDate != 0 {
			includeArtist = artist.CreationDate == ga.creationDate
		}

		if includeArtist {
			if allUnchecked || ga.checkedMembers[len(artist.Members)] {
				card := ga.createCard(artist)
				filtered = append(filtered, card)
			}
		}
	}

	return filtered
}
