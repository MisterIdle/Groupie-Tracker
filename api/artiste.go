package api

type Artist struct {
	Id           int64    `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate uint16   `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsUrl string `json:"locations"`
	ConcertDates string `json:"concertDates"`
	Relations string `json:"relations"`
}
