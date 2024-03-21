package main

import (
	"encoding/json"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

func main() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")
	w.Resize(fyne.NewSize(1000, 600))
	w.SetFixedSize(true)
	w.SetIcon(theme.VolumeUpIcon())

	label := widget.NewLabelWithStyle("Groupie Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	search := widget.NewEntry()
	search.SetPlaceHolder("Search a group or artist")

	artists, err := fetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	suggestionsBox := container.NewVBox()

	header := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewVBox(
			label,
			search,
			suggestionsBox,
		),
	)

	tabs := container.NewAppTabs()

	content := container.New(layout.NewBorderLayout(header, nil, nil, nil),
		header,
		container.NewVScroll(container.NewGridWithColumns(3,
			createArtistCards(artists, tabs)...,
		)),
	)

	hometab := container.NewTabItem("Home", content)
	hometab.Icon = theme.HomeIcon()
	tabs.Append(hometab)

	w.SetContent(tabs)

	search.OnChanged = func(query string) {
		updateSuggestions(query, artists, suggestionsBox)
	}

	w.ShowAndRun()
}

// Suggestions box

func updateSuggestions(query string, artists []Artist, suggestionBox *fyne.Container) {
	suggestionBox.Objects = nil

	if query != "" {
		filtered := filterArtistsAndGroups(query, artists)
		for _, item := range filtered {
			if len(suggestionBox.Objects) >= 6 {
				break
			}
			labelText := item.Name
			if len(item.Members) > 0 {
				labelText += " (" + strings.Join(item.Members, ", ") + ")"
			}
			button := widget.NewButton(labelText, func() {})
			button.Importance = widget.HighImportance
			button.Alignment = widget.ButtonAlignLeading

			suggestionBox.Add(button)
		}
	}
}

// API functions

func fetchArtists() ([]Artist, error) {
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

	return artists, nil
}

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

// Card creation functions

func createArtistCards(artists []Artist, tabs *container.AppTabs) []fyne.CanvasObject {
	var cards []fyne.CanvasObject

	for _, artist := range artists {
		cards = append(cards, createCard(artist, artist.Image, tabs))
	}

	return cards
}

func createCard(artist Artist, imgPath string, tabs *container.AppTabs) fyne.CanvasObject {
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		log.Printf("Error loading image: %v\n", err)
		return nil
	}

	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(230, 230))

	group := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	btn := widget.NewButton("", func() {
		for _, tab := range tabs.Items {
			if tab.Text == artist.Name {
				tabs.Select(tab)
				return
			}
		}

		// If not open, create a new tab for the artist
		artistDetailsTab := createArtistDetailsTab(artist, tabs)
		tabs.Append(container.NewTabItem(artist.Name, artistDetailsTab))
		tabs.Select(tabs.Items[len(tabs.Items)-1])

		// Add a icon to the tab
		tabs.Items[len(tabs.Items)-1].Icon = res

	})

	btn.Importance = widget.LowImportance

	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(1, 30))

	paddedContainer := container.NewPadded(container.New(
		layout.NewStackLayout(),
		btn,
		img,
	))

	vBoxContainer := container.NewVBox(
		space,
		group,
		paddedContainer,
	)

	return vBoxContainer
}

func createArtistDetailsTab(artist Artist, tabs *container.AppTabs) fyne.CanvasObject {
	nameLabel := widget.NewLabel(artist.Name)

	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		// Find the index of the tab to close
		for i, tab := range tabs.Items {
			if tab.Text == artist.Name {
				tabs.RemoveIndex(i)
				return
			}
		}
	})

	nameWithClose := container.NewHBox(
		nameLabel,
		closeButton,
	)

	firstAlbumLabel := widget.NewLabel("First Album: " + artist.FirstAlbum)
	creationDateLabel := widget.NewLabel("Creation Date: " + strconv.Itoa(artist.CreationDate))

	detailsContent := container.NewVBox(
		nameWithClose,
		firstAlbumLabel,
		creationDateLabel,
	)

	detailsScroll := container.NewScroll(detailsContent)

	return detailsScroll
}
