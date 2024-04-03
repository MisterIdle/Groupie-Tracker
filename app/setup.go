// Code by Alexy HOUBLOUP

package app

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GroupieApp struct holds the application window and artists data
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

var (
	countries   []string
	err         error
	maxMembers  = 0
	changeColor = false
	cityMap     = make(map[string]map[string]bool)
)

// Run launches the Groupie application
func (ga *GroupieApp) Run() {
	// Create a new instance of Fyne application
	app := app.New()
	ga.window = app.NewWindow("Groupie Tracker")
	ga.window.Resize(fyne.NewSize(1000, 800))
	ga.window.SetFixedSize(true)
	ga.window.SetIcon(theme.VolumeUpIcon())

	// Create a label for the application title
	label := widget.NewLabelWithStyle("Groupie Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create the search field
	ga.search = widget.NewEntry()
	ga.search.SetPlaceHolder("Search a group or artist")

	// Labels for filters
	memberText := widget.NewLabel("Members:")
	cityLabel := widget.NewLabel("City:")
	sliderLabel := widget.NewLabel("Creation date:")

	// Fetch artists
	ga.artists, err = fetchArtists()
	if err != nil {
		fmt.Println("Error fetching artists:", err)
	}

	// Find the artist with the most members
	for _, artist := range ga.artists {
		if len(artist.Members) > maxMembers {
			maxMembers = len(artist.Members)
		}
	}

	// Create checkboxes for the number of members
	memberCheckboxes := make([]fyne.CanvasObject, maxMembers)
	for i := 0; i < maxMembers; i++ {
		num := i + 1
		memberCheckboxes[i] = widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			ga.checkedMembers[num] = checked
			ga.creationDateToggle.Checked = true
			ga.searchArtists(ga.search.Text)
			ga.cityDropdown.Selected = "All"
		})
	}

	// Create dropdown for cities
	ga.cityDropdown = widget.NewSelect([]string{"All"}, func(city string) {
		ga.searchArtists(ga.search.Text)
	})

	ga.cityDropdown.Selected = "All"

	// Fetch locations for each artist
	for _, artist := range ga.artists {
		artistLocations, err := fetchLocations(artist.ID)
		if err != nil {
			fmt.Println("Error fetching locations:", err)
			continue
		}

		for _, location := range artistLocations {
			parts := strings.Split(location, "-")
			if len(parts) != 2 {
				fmt.Println("Invalid location:", location)
				continue
			}

			country := strings.TrimSpace(strings.ToUpper(parts[1]))
			city := strings.TrimSpace(strings.ReplaceAll(parts[0], "_", " "))
			// Deprecated
			city = strings.Title(city)

			if _, ok := cityMap[country]; !ok {
				cityMap[country] = make(map[string]bool)
			}

			cityMap[country][city] = true
		}
	}

	// Sort countries and cities
	for country := range cityMap {
		countries = append(countries, country)
	}
	sort.Strings(countries)

	// Add countries and cities to the dropdown and sort them
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

	// Update city when dropdown is changed
	ga.cityDropdown.OnChanged = func(selected string) {
		cleanedCity := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(selected, "-", "")))
		cleanedCity = strings.ReplaceAll(cleanedCity, " ", "_")
		ga.city = cleanedCity
		ga.creationDateToggle.Checked = true
		ga.creationDate = 0
		ga.searchArtists(cleanedCity)
	}

	// Fetch min and max creation date
	minDate, maxDate, err := fetchArtistsMinMaxCreationDate()
	if err != nil {
		fmt.Println("Error fetching min and max creation date:", err)
	}

	// Create button to select all dates or specific date
	ga.creationDateToggle = widget.NewCheck("All", func(checked bool) {
		if !checked {
			sliderLabel.SetText(fmt.Sprintf("Creation date: %d", minDate))
			ga.creationDateSlider.SetValue(float64(minDate))
		} else {
			sliderLabel.SetText("Creation date: All")
			ga.creationDate = 0
			ga.searchArtists("")
		}
	})

	ga.creationDateToggle.Checked = true

	// Create slider for creation date
	ga.creationDateSlider = widget.NewSlider(float64(minDate), float64(maxDate))

	ga.creationDateSlider.SetValue(float64(minDate))
	sliderLabel.SetText(fmt.Sprintf("Creation date: %d", minDate))

	// Update creation date when slider is changed
	ga.creationDateSlider.OnChanged = func(value float64) {
		sliderLabel.SetText(fmt.Sprintf("Creation date: %.0f", value))
		ga.creationDateToggle.Checked = false
		ga.creationDate = int(value)
		ga.cityDropdown.Selected = "All"
		valueStr := fmt.Sprintf("%.0f", value)
		ga.searchArtists(valueStr)
	}

	// Create theme button

	// Theme button by Jayanraj rewritten by Alexy HOUBLOUP
	themeButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		changeColor = !changeColor
		if changeColor {
			// Deprecated
			app.Settings().SetTheme(theme.DarkTheme())
		} else {
			// Deprecated
			app.Settings().SetTheme(theme.LightTheme())
		}

		ga.window.SetContent(ga.tabs)
	})

	// UI layout
	ga.suggestionsBox = container.NewVBox()

	membersGroup := container.NewHBox(memberCheckboxes...)

	// Create labels and containers for sliders
	sliderLabelContainer := container.New(layout.NewVBoxLayout(),
		container.NewHBox(
			sliderLabel,
			ga.creationDateToggle,
		),
		ga.creationDateSlider,
	)

	// Create labels and containers for dropdowns
	cityLabelContainer := container.New(layout.NewHBoxLayout(),
		cityLabel,
		ga.cityDropdown,
	)

	// Create filter member container with member checkboxes
	filterMember := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewHBox(
			memberText,
			membersGroup,
		),
	)

	// Create header with search entry, theme button, and filter member container
	header := container.New(layout.NewBorderLayout(nil, nil, nil, nil),
		container.NewVBox(
			container.NewVBox(
				label,
				container.NewBorder(nil, nil, nil, themeButton),
			),
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

	// Set search suggestions
	ga.search.OnChanged = ga.updateSuggestions
	ga.search.OnSubmitted = ga.searchArtists
	ga.checkedMembers = make(map[int]bool)

	// Launch the application
	ga.window.ShowAndRun()
}
