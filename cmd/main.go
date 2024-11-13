package main

import (
	"groupie-tracker/data"
	"groupie-tracker/handlers"
	"log"
	"net/http"
)

// Main function to handle our routes and static files
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/bands", handlers.BandsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch path {
		case "/":
			if r.Method != "GET" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("method not allowed\n"))
				return
			}
			handlers.Homehandler(w, r)
		case "/display":
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("method not allowed\n"))
				return
			}
			handlers.Locationhandler(w, r)
		case "/search":
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("method not allowed\n"))
				return
			}
			handlers.SearchHandler(w, r)
		default:
			data := data.PageData{
				Title: "Error",
				Bands: struct {
					Message string
					Code    int
				}{
					Message: "Page Not Found",
					Code:    404,
				},
			}
			w.WriteHeader(http.StatusNotFound)
			handlers.Rendertemplate(w, data)
		}
	})

	// go-routine to run our server to check if error occurs before running the application
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.Println("running on 8080")

	select {}

}
