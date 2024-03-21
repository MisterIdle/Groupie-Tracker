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

type GroupieApp struct {
	window         fyne.Window
	artists        []Artist
	search         *widget.Entry
	suggestionsBox *fyne.Container
	content        *fyne.Container
	tabs           *container.AppTabs
}

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
	groupieApp := &GroupieApp{}
	groupieApp.Run()
}

func (ga *GroupieApp) Run() {
	a := app.New()
	ga.window = a.NewWindow("Groupie Tracker")
	ga.window.Resize(fyne.NewSize(1000, 800))
	ga.window.SetFixedSize(true)
	ga.window.SetIcon(theme.VolumeUpIcon())

	label := widget.NewLabelWithStyle("Groupie Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	ga.search = widget.NewEntry()
	ga.search.SetPlaceHolder("Search a group or artist")

	var err error
	ga.artists, err = fetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	ga.suggestionsBox = container.NewVBox()

	header := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewVBox(
			label,
			ga.search,
			ga.suggestionsBox,
		),
	)

	ga.tabs = container.NewAppTabs()

	ga.content = container.New(layout.NewBorderLayout(header, nil, nil, nil),
		header,
		container.NewVScroll(container.NewGridWithColumns(3,
			ga.createArtistCards()...,
		)),
	)

	hometab := container.NewTabItem("Home", ga.content)
	hometab.Icon = theme.HomeIcon()
	ga.tabs.Append(hometab)

	ga.window.SetContent(ga.tabs)

	ga.search.OnChanged = ga.updateSuggestions

	ga.search.OnSubmitted = ga.searchArtists

	ga.window.ShowAndRun()
}

func (ga *GroupieApp) searchArtists(query string) {
	filteredCards := ga.filterCards(query)
	if len(filteredCards) == 0 {
		noResultsLabel := widget.NewLabel("Aucun résultat trouvé pour la recherche : " + query)

		// Remplacer le contenu par le message de label
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

	if query != "" {
		filtered := ga.filterArtistsAndGroups(query)
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
			}(item)) // Pass the item (Artist) as an argument to the closure
			button.Importance = widget.HighImportance
			button.Alignment = widget.ButtonAlignLeading

			ga.suggestionsBox.Add(button)
		}
	}
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

func (ga *GroupieApp) createArtistCards() []fyne.CanvasObject {
	var cards []fyne.CanvasObject
	for _, artist := range ga.artists {
		cards = append(cards, ga.createCard(artist, artist.Image))
	}

	return cards
}

func (ga *GroupieApp) createCard(artist Artist, imgPath string) fyne.CanvasObject {
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
		for _, tab := range ga.tabs.Items {
			if tab.Text == artist.Name {
				ga.tabs.Select(tab)
				return
			}
		}

		artistDetailsTab := ga.createArtistDetailsTab(artist)
		ga.tabs.Append(container.NewTabItem(artist.Name, artistDetailsTab))
		ga.tabs.Select(ga.tabs.Items[len(ga.tabs.Items)-1])

		ga.tabs.Items[len(ga.tabs.Items)-1].Icon = res

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

func (ga *GroupieApp) filterCards(query string) []fyne.CanvasObject {
	var filtered []fyne.CanvasObject
	for _, artist := range ga.artists {
		// Vérifier si le nom de l'artiste correspond à la recherche
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			filtered = append(filtered, ga.createCard(artist, artist.Image))
		} else {
			// Vérifier si l'un des membres de l'artiste correspond à la recherche
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
					filtered = append(filtered, ga.createCard(artist, artist.Image))
					break
				}
			}
		}
	}

	return filtered
}

func (ga *GroupieApp) createArtistDetailsTab(artist Artist) fyne.CanvasObject {
	nameLabel := widget.NewLabel(artist.Name)

	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		for i, tab := range ga.tabs.Items {
			if tab.Text == artist.Name {
				ga.tabs.RemoveIndex(i)
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
