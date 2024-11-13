package data

// Structs to capture the json data
type Band struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationdate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Interface to accept any data stucture used in a the bands field
type PageData struct {
	Title string
	Bands interface{}
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Dates_Locations map[string][]string

type Relation struct {
	Id             int             `json:"id"`
	Dateslocations Dates_Locations `json:"datesLocations"`
}


type LocationsApi struct {
	Index []Location `json:"index"`
}