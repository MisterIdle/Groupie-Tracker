package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("Barre de recherche")
	w.Resize(fyne.NewSize(400, 400))

	searchEntry := widget.NewEntry()
	infoLabel := widget.NewLabel("")

	artistImage := canvas.NewImageFromResource(nil)
	artistImage.SetMinSize(fyne.NewSize(300, 200))
	artistImage.FillMode = canvas.ImageFillContain

	searchButton := widget.NewButton("Search", func() {
		artistName := searchEntry.Text
		artists, err := CallAPI(artistName)
		if err != nil {
			fmt.Println("Erreur lors de la recherche d'artistes:", err)
			return
		}

		var foundArtist Artist
		for _, artist := range artists {
			if artist.Name == artistName {
				foundArtist = artist
				break
			}
		}

		if foundArtist.ID != 0 {
			info := printArtistDetails(foundArtist)
			infoLabel.SetText(info)

			img, err := loadImage(foundArtist.Image)
			if err != nil {
				fmt.Println("Erreur lors du chargement de l'image de l'artiste:", err)
				return
			}
			artistImage.Resource = img
		} else {
			infoLabel.SetText("Aucun artiste trouvé.")
			artistImage.Resource = nil
		}

		artistImage.Refresh()
	})

	imageContainer := container.NewCenter(artistImage)
	content := container.NewVBox(
		searchEntry,
		searchButton,
		imageContainer,
		infoLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func loadImage(url string) (fyne.Resource, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return fyne.NewStaticResource("artist_image", data), nil
}

type APIURLs struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relations string `json:"relation"`
}

type Artist struct {
	ID           int             `json:"id"`
	Image        string          `json:"image"`
	Name         string          `json:"name"`
	Members      []string        `json:"members"`
	CreationDate int             `json:"creationDate"`
	FirstAlbum   string          `json:"firstAlbum"`
	Locations    string          `json:"locations"`
	ConcertDates string          `json:"concertDates"`
	Relations    json.RawMessage `json:"relations"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func CallAPI(artistName string) ([]Artist, error) {
	apiURL := "https://groupietrackers.herokuapp.com/api"

	data, err := fetchData(apiURL)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des données: %v", err)
	}

	var apiURLs APIURLs
	err = json.Unmarshal(data, &apiURLs)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'analyse du JSON: %v", err)
	}

	artistsData, err := fetchData(apiURLs.Artists)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des données des artistes: %v", err)
	}

	var artists []Artist
	err = json.Unmarshal(artistsData, &artists)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'analyse du JSON des artistes: %v", err)
	}

	return artists, nil
}

func printArtistDetails(artist Artist) string {
	fmt.Println("--------------")
	fmt.Printf("ID: %d\n", artist.ID)
	fmt.Printf("Nom: %s\n", artist.Name)
	fmt.Printf("Membres: %v\n", artist.Members)
	fmt.Printf("Date de création: %d\n", artist.CreationDate)
	fmt.Printf("Premier album: %s\n", artist.FirstAlbum)
	fmt.Printf("Image: %s\n", artist.Image)
	fmt.Printf("Lieux: %s\n", artist.Locations)
	fmt.Printf("Dates de concert: %s\n", artist.ConcertDates)
	fmt.Printf("Relations: %d\n", artist.Relations)
	info := fmt.Sprintf("ID: %d\nNom: %s\nMembres: %v\nDate de création: %d\nPremier album: %s\n",
		artist.ID, artist.Name, artist.Members, artist.CreationDate, artist.FirstAlbum)
	return info
}

func printLocationDetails(apiURL, locationsURL string) {
	locationData, err := fetchData(locationsURL)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données sur les lieux:", err)
		return
	}

	var locationInfo map[string]interface{}
	err = json.Unmarshal(locationData, &locationInfo)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du JSON des lieux:", err)
		return
	}

	fmt.Println("Informations sur les lieux:")
	fmt.Printf("ID: %v\n", locationInfo["id"])
	fmt.Printf("Lieux: %v\n", locationInfo["locations"])
	fmt.Printf("Dates: %v\n", locationInfo["dates"])
	fmt.Println("--------------")
}

func printDateDetails(apiURL, datesURL string) {
	dateData, err := fetchData(datesURL)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données sur les dates de concert:", err)
		return
	}

	var dateInfo map[string]interface{}
	err = json.Unmarshal(dateData, &dateInfo)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du JSON des dates de concert:", err)
		return
	}

	fmt.Println("Informations sur les dates de concert:")
	fmt.Printf("ID: %v\n", dateInfo["id"])
	fmt.Printf("Dates: %v\n", dateInfo["dates"])
	fmt.Println("--------------")
}

func fetchData(apiURL string) ([]byte, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Réponse avec un code de statut non OK: %d", response.StatusCode)
	}

	return body, nil
}
