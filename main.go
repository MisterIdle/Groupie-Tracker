package main

import (
	"encoding/json"
	"image/color"
	"io"
	"log"
	"net/http"

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
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.SetIcon(theme.VolumeUpIcon())

	label := widget.NewLabelWithStyle("Groupie Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	search := widget.NewEntry()
	search.SetPlaceHolder("Search a group or artist")

	suggestionsBox := container.New(layout.NewCenterLayout(),
		canvas.NewRectangle(color.Black),
		widget.NewLabel("Suggestions"),
	)
	suggestionsBox.Hide()

	artists, err := fetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	header := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewVBox(
			label,
			container.NewVBox(
				search,
			),
			suggestionsBox,
		),
	)

	content := container.New(layout.NewBorderLayout(header, nil, nil, nil),
		header,
		container.NewVScroll(container.NewGridWithColumns(4,
			createArtistCards(artists)...,
		)),
	)

	w.SetContent(content)

	search.OnChanged = func(s string) {
		if s == "" {
			suggestionsBox.Hide()
		} else {
			suggestionsBox.Show()
		}
	}

	w.ShowAndRun()
}

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

func createArtistCards(artists []Artist) []fyne.CanvasObject {
	var cards []fyne.CanvasObject

	for _, artist := range artists {
		cards = append(cards, createCard(artist, artist.Image))
	}

	return cards
}

func createCard(artist Artist, imgPath string) fyne.CanvasObject {
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		log.Printf("Error loading image: %v\n", err)
		return nil
	}

	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))

	btn := widget.NewButton("Open", func() {
		log.Printf("Opening artist: %s\n", artist.Name)
	})

	container := fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		btn,
		img,
	)

	return container
}
