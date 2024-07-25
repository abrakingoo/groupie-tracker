package main

import (
    "html/template"
    "io"
    "log"
    "net/http"
    "encoding/json"
	"strconv"
	"time"
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
    Id    int      `json:"id"`
    Dates []string `json:"dates"`
}

var httpClient = &http.Client{
    Timeout: 10 * time.Second, // Set a 10-second timeout
}


const pageSize = 8

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

func responseAlreadyWritten(w http.ResponseWriter) bool {
    return w.Header().Get("Content-Type") != ""
}

func handleError(w http.ResponseWriter, err error, statusCode int, templateName string) {
    log.Println(err)

	if !responseAlreadyWritten(w) {
        w.WriteHeader(statusCode)
    }
    tpl, tplErr := template.ParseFiles("templates/" + templateName)
    if tplErr != nil {
        log.Println("Error parsing error template file:", tplErr)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if execErr := tpl.Execute(w, nil); execErr != nil {
        log.Println("Error executing error template:", execErr)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    url := r.FormValue("url")
    res, err := http.Get(url)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }
    defer res.Body.Close()

    resBody, err := io.ReadAll(res.Body)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    var dates DateStruct
    err = json.Unmarshal(resBody, &dates)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    tpl, err := template.ParseFiles("templates/dates.html")
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    err = tpl.Execute(w, dates)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
    }
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    url := r.FormValue("url")
    res, err := http.Get(url)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }
    defer res.Body.Close()

    resBody, err := io.ReadAll(res.Body)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    var locations Location
    err = json.Unmarshal(resBody, &locations)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    tpl, err := template.ParseFiles("templates/locations.html")
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    err = tpl.Execute(w, locations)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
    }
}

var templateFuncs = template.FuncMap{
    "sub": func(x, y int) int { return x - y },
    "add": func(x, y int) int { return x + y },
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Get the page query parameter
    page := 1
    if r.URL.Query().Get("page") != "" {
        page, _ = strconv.Atoi(r.URL.Query().Get("page"))
    }

    url := "https://groupietrackers.herokuapp.com/api/artists"
    res, err := http.Get(url)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }
    defer res.Body.Close()

    resBody, err := io.ReadAll(res.Body)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    var bands []Band
    err = json.Unmarshal(resBody, &bands)
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    // Paginate the bands
    paginatedBands, totalItems := Paginate(bands, page)
    totalPages := (totalItems + pageSize - 1) / pageSize

    tpl, err := template.New("index.html").Funcs(templateFuncs).ParseFiles("templates/index.html")
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
        return
    }

    err = tpl.Execute(w, struct {
        Bands       []Band
        CurrentPage int
        TotalPages  int
    }{
        Bands:       paginatedBands,
        CurrentPage: page,
        TotalPages:  totalPages,
    })
    if err != nil {
        handleError(w, err, http.StatusInternalServerError, "500.html")
    }
}


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    tpl, err := template.ParseFiles("templates/404.html")
    if err != nil {
        log.Println("Error parsing 404 template file:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNotFound)
    err = tpl.Execute(w, nil)
    if err != nil {
        log.Println("Error executing 404 template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
            NotFoundHandler(w, r)
        }
    })

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
