package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Band struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type DateStruct struct {
	Id int `json:"id"`
	Dates []string `json:"dates"`
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	url := r.FormValue("url")
	// bandName := r.FormValue("bandName")

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error getting response from API:", err)
		http.Error(w, "Error getting response from API", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var dates DateStruct
	err = json.Unmarshal(resBody, &dates)
	if err != nil {
		log.Println("Error unmarshaling response body:", err)
		http.Error(w, "Error unmarshaling response body", http.StatusInternalServerError)
		return
	}


	tpl, err := template.ParseFiles("templates/dates.html")
	if err != nil {
		log.Println("Error parsing template file:", err)
		http.Error(w, "Error parsing template file", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, dates)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	url := r.FormValue("url")
	// bandName := r.FormValue("bandName")

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error getting response from API:", err)
		http.Error(w, "Error getting response from API", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var locations Location
	err = json.Unmarshal(resBody, &locations)
	if err != nil {
		log.Println("Error unmarshaling response body:", err)
		http.Error(w, "Error unmarshaling response body", http.StatusInternalServerError)
		return
	}


	tpl, err := template.ParseFiles("templates/locations.html")
	if err != nil {
		log.Println("Error parsing template file:", err)
		http.Error(w, "Error parsing template file", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, locations)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error getting response from API:", err)
		http.Error(w, "Error getting response from API", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var bands []Band
	err = json.Unmarshal(resBody, &bands)
	if err != nil {
		log.Println("Error unmarshaling response body:", err)
		http.Error(w, "Error unmarshaling response body", http.StatusInternalServerError)
		return
	}

	tpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template file:", err)
		http.Error(w, "Error parsing template file", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, bands[:8])
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch path {
		case "/":
			HomeHandler(w, r)
		case "/locations":
			LocationsHandler(w, r)
		case "/dates":
			DatesHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
