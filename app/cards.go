// Code by Alexy HOUBLOUP

package app

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	cards []fyne.CanvasObject // A slice to store artist cards
)

// createArtistCards creates cards for each artist
func (ga *GroupieApp) createArtistCards() []fyne.CanvasObject {
	for _, artist := range ga.artists {
		cards = append(cards, ga.createCard(artist)) // Create a card for each artist
	}

	return cards
}

// createCard creates a card for the artist
func (ga *GroupieApp) createCard(artist Artist) fyne.CanvasObject {
	// Load artist image
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return nil
	}

	// Create album cover image and group text
	albumCoverImg := canvas.NewImageFromResource(res)
	albumCoverImg.FillMode = canvas.ImageFillContain
	albumCoverImg.SetMinSize(fyne.NewSize(230, 230))
	groupText := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create button to show artist details
	coverBtn := widget.NewButton("", func() {
		for _, tab := range ga.tabs.Items {
			if tab.Text == artist.Name {
				ga.tabs.Select(tab)
				return
			}
		}

		// If tab for this artist doesn't exist yet, create it
		artistDetailsTab := ga.createArtistDetailsTab(artist)
		ga.tabs.Append(container.NewTabItem(artist.Name, artistDetailsTab))
		ga.tabs.Select(ga.tabs.Items[len(ga.tabs.Items)-1])

		ga.tabs.Items[len(ga.tabs.Items)-1].Icon = res
	})
	coverBtn.Importance = widget.LowImportance

	// Create space between card elements
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(1, 30))

	paddedContainer := container.NewPadded(container.New(
		layout.NewStackLayout(),
		coverBtn,
		albumCoverImg,
	))

	vBoxContainer := container.NewVBox(
		space,
		paddedContainer, groupText,
	)

	return vBoxContainer
}

// createArtistDetailsTab creates a tab for the artist details
func (ga *GroupieApp) createArtistDetailsTab(artist Artist) fyne.CanvasObject {
	// Load artist image
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return nil
	}

	// Create artist image
	cAI := canvas.NewImageFromResource(res)
	cAI.FillMode = canvas.ImageFillContain
	cAI.SetMinSize(fyne.NewSize(230, 230))

	// Create artist details
	nameLabel := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	firstAlbumLabel := widget.NewLabelWithStyle("First Album: "+artist.FirstAlbum, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	creationDateLabel := widget.NewLabelWithStyle("Creation Date: "+strconv.Itoa(artist.CreationDate), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	membreText := widget.NewLabelWithStyle("Members", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	membersLabel := widget.NewLabelWithStyle(strings.Join(artist.Members, "\n"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create close button for tab
	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		for i, tab := range ga.tabs.Items {
			if tab.Text == artist.Name {
				ga.tabs.RemoveIndex(i)
				return
			}
		}
	})

	// Container for artist details with a close button
	detailsContent := container.NewVBox(
		container.NewBorder(nil, nil, nil, closeButton, nameLabel),

		nameLabel,
		cAI,

		membreText,
		membersLabel,

		firstAlbumLabel,
		creationDateLabel,
	)

	// Resize content to fit minimum size
	detailsContent.Resize(detailsContent.MinSize())

	// Create scroll container for artist details
	detailsScroll := container.NewScroll(detailsContent)
	return detailsScroll
}
