package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"groupie-tracker/handlers"
)

func TestRouterAndHomehandler(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"GET Root", http.MethodGet, "/", http.StatusOK, ""},                                   
		{"POST Root", http.MethodPost, "/", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"PUT Root", http.MethodPut, "/", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"GET Other Path", http.MethodGet, "/other", http.StatusNotFound, "404 page not found\n"}, 
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method and path
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a handler that mimics the routing logic
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				switch path {
				case "/":
					if r.Method != "GET" {
						w.WriteHeader(http.StatusMethodNotAllowed)
						w.Write([]byte("method not allowed\n"))
						return
					}
					handlers.Homehandler(w, r)
				default:
					http.NotFound(w, r)
				}
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			body, _ := io.ReadAll(rr.Body)
			if tc.expectedBody != "" && strings.TrimSpace(string(body)) != strings.TrimSpace(tc.expectedBody) {
				t.Errorf("Handler returned unexpected body: got %v want %v",
					string(body), tc.expectedBody)
			}
		})
	}
}



func TestRouterAndLocationHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{                                   
		{"POST Root", http.MethodGet, "/display", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"PUT Root", http.MethodPut, "/display", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"GET Other Path", http.MethodGet, "/other", http.StatusNotFound, "404 page not found\n"}, 
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method and path
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a handler that mimics the routing logic
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				switch path {
				case "/display":
					if r.Method != "POST" {
						w.WriteHeader(http.StatusMethodNotAllowed)
						w.Write([]byte("method not allowed\n"))
						return
					}
					handlers.Locationhandler(w, r)
				default:
					http.NotFound(w, r)
				}
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			body, _ := io.ReadAll(rr.Body)
			if tc.expectedBody != "" && strings.TrimSpace(string(body)) != strings.TrimSpace(tc.expectedBody) {
				t.Errorf("Handler returned unexpected body: got %v want %v",
					string(body), tc.expectedBody)
			}
		})
	}
}

func TestRouterAndSearchHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{                                  
		{"POST Root", http.MethodGet, "/search", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"PUT Root", http.MethodPut, "/search", http.StatusMethodNotAllowed, "method not allowed\n"},
		{"GET Other Path", http.MethodGet, "/other", http.StatusNotFound, "404 page not found\n"}, 
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method and path
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a handler that mimics the routing logic
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				switch path {
				case "/search":
					if r.Method != "POST" {
						w.WriteHeader(http.StatusMethodNotAllowed)
						w.Write([]byte("method not allowed\n"))
						return
					}
					handlers.SearchHandler(w, r)
				default:
					http.NotFound(w, r)
				}
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			body, _ := io.ReadAll(rr.Body)
			if tc.expectedBody != "" && strings.TrimSpace(string(body)) != strings.TrimSpace(tc.expectedBody) {
				t.Errorf("Handler returned unexpected body: got %v want %v",
					string(body), tc.expectedBody)
			}
		})
	}
}