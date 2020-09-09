package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

// TmplMap maps html template paths to shortnames
var TmplMap = map[string]string{
	"Index":          "../../frontend/html/Index.html",
	"rankings":       "../../frontend/html/rankings.html",
	"about":          "../../frontend/html/about.html",
	"search-results": "../../frontend/html/search-results.html",
}

// InitHTTPServer initializes an http server at the provided address
func InitHTTPServer(addr string) *http.Server {
	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
		Handler:      nil,
	}
	return srv
}

// RegisterHandlers registers the http handler functions:
//  Home
//  Rankings
//  About
//  SearchResults
// and creates file servers for the css, js, and img files.
func RegisterHandlers() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/rankings", Rankings)
	http.HandleFunc("/about", About)
	http.HandleFunc("/search-results/", SearchResults)
	if _, err := os.Stat("../../frontend/css"); os.IsNotExist(err) {
		fmt.Printf("WARNING: css files not found at '../../frontend/css'\n")
	}
	cssHandler := http.StripPrefix("/css/", http.FileServer(http.Dir("../../frontend/css")))
	http.Handle("/css/", cssHandler)
	if _, err := os.Stat("../../frontend/js"); os.IsNotExist(err) {
		fmt.Printf("WARNING: css files not found at '../../frontend/js'\n")
	}
	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("../../frontend/js")))
	http.Handle("/js/", jsHandler)
}

// Home displays home page
func Home(w http.ResponseWriter, r *http.Request) {
	if err := templExe(TmplMap["Index"], w, nil); err != nil {
		// util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	return
}

// Rankings displays rankings page
func Rankings(w http.ResponseWriter, r *http.Request) {
	if err := templExe(TmplMap["rankings"], w, nil); err != nil {
		// util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	return
}

// About displays About page
func About(w http.ResponseWriter, r *http.Request) {
	if err := templExe(TmplMap["about"], w, nil); err != nil {
		// util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	return
}

// SearchResults displays the search results page
func SearchResults(w http.ResponseWriter, r *http.Request) {
	if err := templExe(TmplMap["search-results"], w, nil); err != nil {
		// util.FailLog(err)
		fmt.Fprintf(w, "html template failed to execute: %s", err)
		fmt.Printf("html template failed to execute: %s", err)
		return
	}
	return
}

// templExe executes for the specified template for the given io.Writer and data interface
func templExe(tmpl string, w http.ResponseWriter, data interface{}) error {
	t := template.Must(template.ParseFiles(tmpl))
	if err := t.Execute(w, data); err != nil {
		fmt.Printf("template execution failed: %v", err)
		return fmt.Errorf("template execution failed: %v", err)
	}
	return nil
}
