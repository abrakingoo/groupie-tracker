package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Band struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relations struct {
    Id             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}

func FetchData(url string) ([]Band, error) {
	var Bandslice []Band

	res, err := http.Get(url)
	if err != nil {
		return Bandslice, fmt.Errorf("error fetching data: %v", err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Bandslice, err
	}

	err = json.Unmarshal(resBody, &Bandslice)
	if err != nil {
		return Bandslice, err
	}

	return Bandslice, nil
}

func processURLHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()

	// Get the URL parameter from the form data
	url := r.FormValue("url")

	// Fetch the data from the URL
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Read the response body
	resBody, _ := io.ReadAll(res.Body)

	log.Printf("RES BODY: %v", string(resBody))
	// Define the Relations struct
	var relations Relations

	// Unmarshal JSON into the Relations struct
	err = json.Unmarshal(resBody, &relations)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	
	// Parse the HTML template
	tmp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the Relations data
	err = tmp.Execute(w, relations)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

type Res struct {
	Data []Band
}

func main() {
	artists := "https://groupietrackers.herokuapp.com/api/artists"
	// locations := "https://groupietrackers.herokuapp.com/api/locations"
	// dates := "https://groupietrackers.herokuapp.com/api/dates"
	// relation := "https://groupietrackers.herokuapp.com/api/relation"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/process-url", processURLHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
		}

		res, err := FetchData(artists)
		// res, err := FetchData(relation)
		if err != nil {
			return
		}

		resdata := Res{Data: res}

		tmpl.Execute(w, resdata)
	})

	log.Println("Server Started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
