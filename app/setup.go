package app

import (
	"fmt"
	"log"

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
	checkedMembers map[int]bool // Stocke l'état des cases cochées
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

	memberCheckboxes := make([]fyne.CanvasObject, maxMembers)
	for i := 0; i < maxMembers; i++ {
		num := i + 1
		memberCheckboxes[i] = widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			ga.checkedMembers[num] = checked
			ga.searchArtists(ga.search.Text)
		})
	}

	ga.search = widget.NewEntry()
	ga.search.SetPlaceHolder("Search a group or artist")

	ga.suggestionsBox = container.NewVBox()

	membersGroup := container.NewHBox(memberCheckboxes...)

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
}
