// Code by Alexy HOUBLOUP

package app

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// searchArtists filters artists based on the search query and updates the display
func (ga *GroupieApp) searchArtists(query string) {
	if query == "all" {
		query = ""
	}

	// Reset city filter if query is empty
	if query == "" {
		ga.cityDropdown.Selected = "All"
	}

	// Filter artist cards
	filteredCards := ga.filterCards(query)

	// If no results found, display appropriate message
	if len(filteredCards) == 0 {
		noResultsLabel := widget.NewLabel("No results found for search: " + query)
		ga.content.Objects[1] = noResultsLabel
		ga.content.Refresh()
		ga.search.SetText("")
		return
	}

	// Update content with filtered cards
	filteredContent := container.NewVScroll(container.NewGridWithColumns(3, filteredCards...))
	ga.content.Objects[1] = filteredContent
	ga.content.Refresh()
	ga.search.SetText("")
}

// updateSuggestions updates the suggestions box based on the search query
func (ga *GroupieApp) updateSuggestions(query string) {
	ga.suggestionsBox.Objects = nil
	var filtered []Artist

	if query != "" {
		// Convert query to integer (if numeric)
		queryInt, err := strconv.Atoi(query)

		// Filter artists based on query
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

			// Filter by creation date if query is a number
			if err == nil && artist.CreationDate == queryInt {
				filtered = append(filtered, artist)
			}

			if strings.Contains(artist.FirstAlbum, query) {
				filtered = append(filtered, artist)
			}
		}

		// Other filters based on city and members
		citySelected := ga.city != "all"
		allUnchecked := true
		for _, checked := range ga.checkedMembers {
			if checked {
				allUnchecked = false
				break
			}
		}

		// Add suggestions based on filters
		for _, item := range filtered {
			if len(ga.suggestionsBox.Objects) >= 5 {
				break
			}

			if allUnchecked || ga.checkedMembers[len(item.Members)] || !citySelected {
				loc, err := fetchLocations(item.ID)
				if err != nil {
					fmt.Println("Error fetching locations for artist", item.Name, ":", err)
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

// containsLocation checks if the artist or group has a location that contains the city
func containsLocation(locations []string, city string) bool {
	for _, loc := range locations {
		if strings.Contains(strings.ToLower(loc), strings.ToLower(city)) {
			return true
		}
	}
	return false
}

// filterArtistByLocation filters artists based on the location
func (ga *GroupieApp) filterArtistByLocation(location string) []Artist {
	var filtered []Artist

	for _, artist := range ga.artists {
		locations, err := fetchLocations(artist.ID)
		if err != nil {
			fmt.Println("Error fetching locations for artist", artist.Name, ":", err)
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

// filterCards filters cards based on the query
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

		if dateSearch && artist.CreationDate == queryInt {
			includeArtist = true
		}

		if query != "" {
			locations := ga.filterArtistByLocation(query)
			for _, locArtist := range locations {
				if locArtist.ID == artist.ID {
					includeArtist = true
					break
				}
			}
		}

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
