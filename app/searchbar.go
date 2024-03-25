package app

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ga *GroupieApp) searchArtists(query string) {
	ga.showLoadingText(true)

	filteredCards := ga.filterCards(query)

	ga.showLoadingText(false)

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

func (ga *GroupieApp) showLoadingText(show bool) {
	if show {
		ga.content.Objects[1] = widget.NewLabel("Chargement...")
	} else {
		ga.content.Objects[1] = container.NewVScroll(container.NewGridWithColumns(3, ga.createArtistCards()...))
	}
}

func (ga *GroupieApp) updateSuggestions(query string) {
	ga.suggestionsBox.Objects = nil

	if query != "" {
		var filtered []Artist

		queryInt, err := strconv.Atoi(query)
		if err == nil {
			filtered = ga.filterArtistsByCreationDate(queryInt)
		} else if strings.Contains(query, "") {
			filtered = ga.filterArtistByFirstAlbum(query)
		} else {
			filtered = ga.filterArtistsAndGroups(query)
		}

		for _, item := range filtered {
			if len(ga.suggestionsBox.Objects) >= 6 {
				break
			}
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

func (ga *GroupieApp) filterArtistsAndGroups(query string) []Artist {
	var filtered []Artist

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
	}

	return filtered
}

func (ga *GroupieApp) filterArtistsByCreationDate(date int) []Artist {
	var filtered []Artist

	for _, artist := range ga.artists {
		yearStr := strconv.Itoa(artist.CreationDate)

		// Si la date est contenue dans l'année de création de l'artiste
		if strings.Contains(yearStr, strconv.Itoa(date)) {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

// "firstAlbum":"14-12-1973"
func (ga *GroupieApp) filterArtistByFirstAlbum(album string) []Artist {
	var filtered []Artist

	for _, artist := range ga.artists {
		if strings.Contains(artist.FirstAlbum, album) {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}

func (ga *GroupieApp) filterCards(query string) []fyne.CanvasObject {
	queryLower := strings.ToLower(query)
	var filtered []fyne.CanvasObject

	queryInt, err := strconv.Atoi(query)
	var dateSearch bool

	if err == nil {
		dateSearch = true
	}

	// Si la requête est vide, retourner simplement toutes les cartes
	if query == "" {
		return ga.createArtistCards()
	}

	for _, artist := range ga.artists {
		if strings.Contains(strings.ToLower(artist.Name), queryLower) {
			filtered = append(filtered, ga.createCard(artist))
		} else {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), queryLower) {
					filtered = append(filtered, ga.createCard(artist))
					break
				}
			}
		}
		if dateSearch && artist.CreationDate == queryInt {
			filtered = append(filtered, ga.createCard(artist))
		}

		if strings.Contains(artist.FirstAlbum, query) {
			filtered = append(filtered, ga.createCard(artist))
		}
	}

	return filtered
}
