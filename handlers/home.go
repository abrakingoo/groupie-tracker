package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"groupie-tracker/data"
)

var (
	tpl *template.Template
	err error
)

// Render the temmplate before anything else
func init() {
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		return
	}
}

// Render the template with data
func Rendertemplate(w http.ResponseWriter, data interface{}) {
	if tpl == nil {
		return
	}
	tpl.ExecuteTemplate(w, "base.html", data)
}

var Bandis []data.Band

// Fetch and serve our json data to the homepage
func Homehandler(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Error fetching data",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		Rendertemplate(w, data)
		return
	}

	defer resp.Body.Close()
	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Unexpected status code from external service",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		Rendertemplate(w, data)
		return

	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Error reading response data",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		Rendertemplate(w, data)
		return
	}


	err = json.Unmarshal(resBody, &Bandis)
	if err != nil {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Error Unmarshaling data",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		Rendertemplate(w, data)
		return
	}

	data := data.PageData{
		Title: "Home",
		Bands: Bandis,
	}
	Rendertemplate(w, data)
}
