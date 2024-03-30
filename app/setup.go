package app

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	checkedMembers map[int]bool
}

func (ga *GroupieApp) Run() {
	a := app.New()
	ga.window = a.NewWindow("Groupie Tracker")
	ga.window.Resize(fyne.NewSize(1000, 800))
	ga.window.SetFixedSize(true)
	ga.window.SetIcon(theme.VolumeUpIcon())

	label := widget.NewLabelWithStyle("Groupie Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	var err error
	ga.artists, err = fetchArtists()
	if err != nil {
		log.Fatal("Error fetching artists:", err)
	}

	maxMembers := 0
	for _, artist := range ga.artists {
		if len(artist.Members) > maxMembers {
			maxMembers = len(artist.Members)
		}
	}

	memberText := widget.NewLabel("Members:")
	cityLabel := widget.NewLabel("City:")

	memberCheckboxes := make([]fyne.CanvasObject, maxMembers)
	for i := 0; i < maxMembers; i++ {
		num := i + 1
		memberCheckboxes[i] = widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			ga.checkedMembers[num] = checked
			ga.searchArtists(ga.search.Text)

			fmt.Println(ga.checkedMembers)
		})
	}

	cityDropdown := widget.NewSelect([]string{"All"}, func(city string) {
		ga.searchArtists(ga.search.Text)
	})

	cityMap := make(map[string]map[string]bool)

	for _, artist := range ga.artists {
		artistLocations, err := fetchLocations(artist.ID)
		if err != nil {
			log.Printf("Error fetching locations for artist %s: %v\n", artist.Name, err)
			continue
		}

		for _, location := range artistLocations {
			parts := strings.Split(location, "-")
			if len(parts) != 2 {
				log.Printf("Invalid location format: %s\n", location)
				continue
			}

			country := strings.TrimSpace(strings.ToUpper(parts[1]))
			city := strings.TrimSpace(strings.ReplaceAll(parts[0], "_", " "))
			city = strings.Title(city)

			if _, ok := cityMap[country]; !ok {
				cityMap[country] = make(map[string]bool)
			}

			cityMap[country][city] = true
		}
	}

	// Tri par pays
	var countries []string
	for country := range cityMap {
		countries = append(countries, country)
	}
	sort.Strings(countries)

	for _, country := range countries {
		cityDropdown.Options = append(cityDropdown.Options, country)

		cities := make([]string, 0)
		for city := range cityMap[country] {
			cities = append(cities, city)
		}
		sort.Strings(cities)

		for _, city := range cities {
			cityDropdown.Options = append(cityDropdown.Options, fmt.Sprintf("  - %s", city))
		}
	}

	ga.search = widget.NewEntry()
	ga.search.SetPlaceHolder("Search a group or artist")

	ga.suggestionsBox = container.NewVBox()

	membersGroup := container.NewHBox(memberCheckboxes...)

	cityLabelContainer := container.New(layout.NewHBoxLayout(),
		cityLabel,
		cityDropdown,
	)

	filterMember := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewHBox(
			memberText,
			membersGroup,
		),
	)

	header := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewVBox(
			label,
			filterMember,
			cityLabelContainer,
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

	ga.checkedMembers = make(map[int]bool) // Initialiser checkedMembers

	ga.window.ShowAndRun()

	//code poour le bouton pour changer la couleur de l'app
	w := a.NewWindow("GROUPIE-TRACKER")
	// Créer un bouton pour changer la couleur de l'application
	var changeColor bool
	button := widget.NewButton("Change the background color", func() {
		changeColor = !changeColor
		if changeColor {
			// Changer la couleur de l'application en noir et blanc
			theme := theme.DarkTheme()
			fyne.CurrentApp().Settings().SetTheme(theme)
		} else {
			// Rétablir la couleur de l'application par défaut
			theme := theme.LightTheme()
			fyne.CurrentApp().Settings().SetTheme(theme)
		}
		// Appliquer les changements de couleur à la fenêtre
		w.Canvas().Refresh(w.Content())
	})

	// Créer un conteneur pour centrer le bouton
	content := container.NewCenter(button)

	// Définir le contenu de la fenêtre
	w.SetContent(content)
}
