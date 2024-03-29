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
	window             fyne.Window
	artists            []Artist
	search             *widget.Entry
	suggestionsBox     *fyne.Container
	content            *fyne.Container
	tabs               *container.AppTabs
	checkedMembers     map[int]bool
	cityDropdown       *widget.Select
	city               string
	creationDateSlider *widget.Slider
	creationDateToggle *widget.Check
	creationDate       int
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
	sliderLabel := widget.NewLabel("Creation date:")

	memberCheckboxes := make([]fyne.CanvasObject, maxMembers)
	for i := 0; i < maxMembers; i++ {
		num := i + 1
		memberCheckboxes[i] = widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			ga.checkedMembers[num] = checked
			ga.searchArtists(ga.search.Text)
			ga.cityDropdown.Selected = "All"
		})
	}

	ga.cityDropdown = widget.NewSelect([]string{"All"}, func(city string) {
		ga.searchArtists(ga.search.Text)
	})

	ga.cityDropdown.Selected = "All"

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

	var countries []string
	for country := range cityMap {
		countries = append(countries, country)
	}
	sort.Strings(countries)

	for _, country := range countries {
		ga.cityDropdown.Options = append(ga.cityDropdown.Options, country)

		cities := make([]string, 0)
		for city := range cityMap[country] {
			cities = append(cities, city)
		}
		sort.Strings(cities)

		for _, city := range cities {
			ga.cityDropdown.Options = append(ga.cityDropdown.Options, fmt.Sprintf("  - %s", city))
		}
	}

	ga.cityDropdown.OnChanged = func(selected string) {
		cleanedCity := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(selected, "-", "")))
		cleanedCity = strings.ReplaceAll(cleanedCity, " ", "_")
		ga.city = cleanedCity
		ga.searchArtists(cleanedCity)
	}

	minDate, maxDate, err := fetchArtistsMinMaxCreationDate()
	if err != nil {
		log.Fatal("Error fetching min and max creation date:", err)
	}

	ga.creationDateToggle = widget.NewCheck("All", func(checked bool) {
		if !checked {
			sliderLabel.SetText(fmt.Sprintf("Creation date: %d", minDate))
			ga.creationDateSlider.SetValue(float64(minDate))
			ga.creationDateSlider.Hidden = false
		} else {
			ga.creationDateSlider.Hidden = true
			sliderLabel.SetText("Creation date: All")
			ga.creationDate = 0
			ga.searchArtists("")
		}
	})

	ga.creationDateSlider = widget.NewSlider(float64(minDate), float64(maxDate))

	ga.creationDateSlider.SetValue(float64(minDate))
	sliderLabel.SetText(fmt.Sprintf("Creation date: %d", minDate))

	ga.creationDateSlider.OnChanged = func(value float64) {
		sliderLabel.SetText(fmt.Sprintf("Creation date: %.0f", value))
		ga.creationDateToggle.Checked = false
		ga.creationDate = int(value)
		valueStr := fmt.Sprintf("%.0f", value)
		ga.searchArtists(valueStr)
	}

	ga.search = widget.NewEntry()
	ga.search.SetPlaceHolder("Search a group or artist")

	ga.suggestionsBox = container.NewVBox()

	membersGroup := container.NewHBox(memberCheckboxes...)

	sliderLabelContainer := container.New(layout.NewVBoxLayout(),
		container.NewHBox(
			sliderLabel,
			ga.creationDateToggle,
		),
		ga.creationDateSlider,
	)

	cityLabelContainer := container.New(layout.NewHBoxLayout(),
		cityLabel,
		ga.cityDropdown,
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
			sliderLabelContainer,
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

	ga.checkedMembers = make(map[int]bool)

	ga.window.ShowAndRun()
}
