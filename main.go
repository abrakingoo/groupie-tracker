package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Define your structs
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
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Define page size
const pageSize = 8

// Pagination function
func Paginate(items []Band, page int) ([]Band, int) {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []Band{}, len(items)
	}

	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], len(items)
}

// Template functions
var templateFuncs = template.FuncMap{
	"sub": func(x, y int) int { return x - y },
	"add": func(x, y int) int { return x + y },
}

// Load templates globally
var templates = template.Must(template.New("").Funcs(templateFuncs).ParseGlob("templates/*.html"))

// DatesHandler
func DatesHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")

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

	err = templates.ExecuteTemplate(w, "dates.html", dates)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// LocationsHandler
func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")

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

	err = templates.ExecuteTemplate(w, "locations.html", locations)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// HomeHandler
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

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if pNum, err := strconv.Atoi(p); err == nil {
			page = pNum
		}
	}

	paginatedBands, totalItems := Paginate(bands, page)
	totalPages := (totalItems + pageSize - 1) / pageSize

	data := struct {
		Bands      []Band
		CurrentPage int
		TotalPages int
	}{
		Bands:       paginatedBands,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// main function
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
