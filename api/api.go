package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type APIURLs struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relations string `json:"relation"`
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
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func CallAPI() {
	apiURL := "https://groupietrackers.herokuapp.com/api"

	data, err := fetchData(apiURL)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données:", err)
		return
	}

	var apiURLs APIURLs
	err = json.Unmarshal(data, &apiURLs)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du JSON:", err)
		return
	}

	artistsData, err := fetchData(apiURLs.Artists)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données des artistes:", err)
		return
	}

	var artists []Artist
	err = json.Unmarshal(artistsData, &artists)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du JSON des artistes:", err)
		return
	}

	relationsData, err := fetchData(apiURLs.Relations)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données des relations:", err)
		return
	}

	var relation Relation
	err = json.Unmarshal(relationsData, &relation)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du JSON des relations:", err)
		return
	}

	var inputID int
	fmt.Print("Veuillez saisir l'ID du groupe : ")
	_, err = fmt.Scan(&inputID)
	if err != nil {
		fmt.Println("Erreur lors de la saisie de l'ID du groupe:", err)
		return
	}

	found := false
	for _, artist := range artists {
		if artist.ID == inputID {
			printArtistDetails(apiURLs, artist)
			printLocationDetails(apiURLs.Locations, artist.Locations)
			printDateDetails(apiURLs.Dates, artist.ConcertDates)

			found = true
			break
		}
	}

	if !found {
		fmt.Println("Aucun artiste trouvé avec l'ID spécifié.")
	}
}

func printArtistDetails(apiURLs APIURLs, artist Artist) {
	fmt.Println("--------------")
	fmt.Printf("ID: %d\n", artist.ID)
	fmt.Printf("Nom: %s\n", artist.Name)
	fmt.Printf("Membres: %v\n", artist.Members)
	fmt.Printf("Date de création: %d\n", artist.CreationDate)
	fmt.Printf("Premier album: %s\n", artist.FirstAlbum)
	fmt.Printf("Image: %s\n", artist.Image)
	fmt.Printf("Lieux: %s\n", artist.Locations)
	fmt.Printf("Dates de concert: %s\n", artist.ConcertDates)
	fmt.Printf("Relations: %s\n", apiURLs.Relations+"/"+strconv.Itoa(artist.ID))
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
