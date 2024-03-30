package app

import (
	"image/color"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (ga *GroupieApp) createArtistCards() []fyne.CanvasObject {
	var cards []fyne.CanvasObject
	for _, artist := range ga.artists {
		cards = append(cards, ga.createCard(artist))
	}

	return cards
}

func (ga *GroupieApp) createCard(artist Artist) fyne.CanvasObject {
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
		paddedContainer, group,
	)

	return vBoxContainer
}

func (ga *GroupieApp) createArtistDetailsTab(artist Artist) fyne.CanvasObject {

	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		log.Printf("Error loading image: %v\n", err)
		return nil
	}

	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(230, 230))

	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	firstAlbumLabel := widget.NewLabelWithStyle("First Album: "+artist.FirstAlbum, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	creationDateLabel := widget.NewLabelWithStyle("Creation Date: "+strconv.Itoa(artist.CreationDate), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	membreText := widget.NewLabelWithStyle("Members", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	membersLabel := widget.NewLabelWithStyle(strings.Join(artist.Members, "\n"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		for i, tab := range ga.tabs.Items {
			if tab.Text == artist.Name {
				ga.tabs.RemoveIndex(i)
				return
			}
		}
	})

	detailsContent := container.NewVBox(
		container.NewBorder(nil, nil, nil, closeButton, nameLabel),

		nameLabel,
		img,

		membreText,
		membersLabel,

		firstAlbumLabel,
		creationDateLabel,
	)

	detailsContent.Resize(detailsContent.MinSize())

	detailsScroll := container.NewScroll(detailsContent)
	return detailsScroll
}
