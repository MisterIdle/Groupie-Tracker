package main

import (
	"encoding/json"
	"fmt"
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

type TappableCard struct {
	*fyne.Container
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

	fmt.Println(artists)

	return artists, nil
}

func createArtistCards(artists []Artist) []fyne.CanvasObject {
	var cards []fyne.CanvasObject

	for _, artist := range artists {
		card := createCard(artist)
		if card != nil {
			cards = append(cards, card)
		}
	}

	return cards
}

func createCard(artist Artist) fyne.CanvasObject {
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		log.Printf("Error loading image: %v\n", err)
		return nil
	}

	image := canvas.NewImageFromResource(res)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(200, 200))

	name := widget.NewLabel(artist.Name)
	name.Alignment = fyne.TextAlignCenter
	name.TextStyle = fyne.TextStyle{Bold: true}

	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	border.StrokeWidth = 3

	// Conteneur pour la carte avec contour et contenu
	card := container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(
			container.NewCenter(image),
			container.New(layout.NewVBoxLayout(),
				container.NewCenter(name),
			),
		),
		border,
	)

	return card
}
