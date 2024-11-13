package handlers

import (
	"encoding/json"
	"groupie-tracker/data"
	"io"
	"net/http"
)

// Locationhandler processes form submissions and fetches data from a URL.
func Locationhandler(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	if err := r.ParseForm(); err != nil {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Error Parsing Form",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		Rendertemplate(w, data)
		return
	}

	formType := r.FormValue("form_type")
	url := r.FormValue("url")
	image := r.FormValue("image")
	name := r.FormValue("bandName")

	// Fetch data from the URL
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

	// Initialize data structures
	var location data.Location
	var dates data.Date
	var relations data.Relation

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

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "Error reading response body",
				Code:    500,
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		Rendertemplate(w, data)
		return
	}
	// Switch case to select page and render page data
	switch formType {
	case "Dates":
		if err := json.Unmarshal(body, &dates); err != nil {
			http.Error(w, "Error unmarshalling dates data", http.StatusInternalServerError)
			return
		}
		data := data.PageData{
			Title: "Dates",
			Bands: struct {
				Dates data.Date
				Image string
				Name  string
			}{
				Dates: dates,
				Image: image,
				Name:  name,
			},
		}
		Rendertemplate(w, data)

	case "Location":
		if err := json.Unmarshal(body, &location); err != nil {
			data := data.PageData{
				Title: "Error",
				Bands: "Error unmarshalling location data",
			}
			w.WriteHeader(http.StatusInternalServerError)
			Rendertemplate(w, data)
		}
		data := data.PageData{
			Title: "Location",
			Bands: struct {
				Locations data.Location
				Image     string
				Name      string
			}{
				Locations: location,
				Image:     image,
				Name:      name,
			},
		}
		Rendertemplate(w, data)

	case "Relations":
		if err := json.Unmarshal(body, &relations); err != nil {
			http.Error(w, "Error unmarshalling relations data", http.StatusInternalServerError)
			return
		}
		data := data.PageData{
			Title: "Relations",
			Bands: struct {
				Relations data.Relation
				Image     string
				Name      string
			}{
				Relations: relations,
				Image:     image,
				Name:      name,
			},
		}
		Rendertemplate(w, data)

	default:
		data := data.PageData{
			Title: "Error",
			Bands: struct {
				Message string
				Code    int
			}{
				Message: "unknown form type",
				Code:    500,
			},
		}
		Rendertemplate(w, data)
	}
}
