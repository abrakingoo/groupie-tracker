package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/data"
)

// Utility function to handle error responses
func handleError(w http.ResponseWriter, statusCode int, message string) {
	data := data.PageData{
		Title: "Error", // Title for the error page
		Bands: struct {
			Message string
			Code    int
		}{
			Message: message, // Error message to display
			Code:    statusCode, // HTTP status code
		},
	}
	w.WriteHeader(statusCode) // Set the HTTP status code
	Rendertemplate(w, data) // Render the error template
}

// Check if a specific band is already present in the list of bands
func isBand(bands []data.Band, band data.Band) bool {
	for _, bnd := range bands {
		if bnd.Name == band.Name { // Compare band names
			return true // Band found
		}
	}
	return false // Band not found
}

// SearchHandler handles the incoming search requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var band []data.Band // Slice to hold matched bands
	if r.Method != http.MethodPost { // Check for POST method
		handleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Parse form data from the request
	if err := r.ParseForm(); err != nil {
		handleError(w, http.StatusBadRequest, "Error Parsing Form")
		return
	}

	// Split the search query by " - "
	query := strings.Split(r.FormValue("search"), " - ")
	if len(query) == 0 {
		handleError(w, http.StatusBadRequest, "invalid search term")
		return
	}
	searchTerm := query[0] // Get the first part of the search term

	// Fetch location data from the external API
	var locationstruct data.LocationsApi
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error fetching data")
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		handleError(w, http.StatusInternalServerError, "Unexpected status code from external service")
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error reading response body")
		return
	}

	// Unmarshal the JSON response into the LocationsApi struct
	if err := json.Unmarshal(body, &locationstruct); err != nil {
		handleError(w, http.StatusInternalServerError, "Error unmarshalling location data")
		return
	}

	// Iterate through all available bands to find matches
	for _, b := range Bandis {
		// Check if the search term is a valid integer and matches the creation date
		if val, err := strconv.Atoi(searchTerm); err == nil {
			if b.CreationDate == val {
				band = append(band, b) // Add band if creation date matches
			}
		}

		// Check if the band's name contains the search term
		if strings.Contains(strings.TrimSpace(strings.ToLower(b.Name)), strings.TrimSpace(strings.ToLower(searchTerm))) {
			band = append(band, b) // Add band if name matches
		}

		// Check if the first album contains the search term
		if strings.Contains(b.FirstAlbum, searchTerm) {
			band = append(band, b) // Add band if first album matches
		}

		// Check if any band members' names contain the search term
		for _, name := range b.Members {
			if strings.Contains(strings.ToLower(name), strings.ToLower(searchTerm)) {
				if !isBand(band, b) { // Avoid adding duplicates
					band = append(band, b)
				}
			}
		}

		// Create a set to track unique location IDs that match the search term
		locationIDs := make(map[int]struct{})
		for _, index := range locationstruct.Index {
			for _, loc := range index.Locations {
				// Check if the location contains the search term
				if strings.Contains(strings.ReplaceAll(strings.ToLower(loc), "_", " "), strings.ToLower(searchTerm)) {
					locationIDs[index.Id] = struct{}{} // Add location ID to the set
				}
			}
		}

		// Add bands that correspond to the matched location IDs
		for id := range locationIDs {
			if id > 0 && id <= len(Bandis) { // Ensure the ID is valid
				if !isBand(band, Bandis[id-1]) {
					band = append(band, Bandis[id-1]) // Add band if not already included
				}
			}
		}
	}

	// Prepare data for rendering the search results
	data := data.PageData{
		Title: "Search", // Title for the search results page
		Bands: struct {
			Band   []data.Band // Slice of matched bands
			Length int         // Length of the matched bands slice
		}{
			Band:   band,
			Length: len(band),
		},
	}
	Rendertemplate(w, data) // Render the results template
}

