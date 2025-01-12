package main

import (
	"bufio"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var apiKey string
var Genres = make(map[int]string)

func init() {
	loadApiKey()
	loadGenres()
}

func loadApiKey() {
	_, currentFile, _, _ := runtime.Caller(1)
	apiKeyFile := path.Join(path.Dir(currentFile), "apiKey")
	file, err := os.Open(apiKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey = scanner.Text()
}

func loadGenres() {
	type Genre struct {
		Id   int
		Name string
	}

	u, _ := url.Parse("https://api.themoviedb.org/3/genre/movie/list")
	q := u.Query()
	q.Add("api_key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode != 200 {
		log.Println("Could not load the genres! " + resp.Status)
	}
	defer resp.Body.Close()

	var result map[string][]Genre
	json.NewDecoder(resp.Body).Decode(&result)

	for _, element := range result["genres"] {
		Genres[element.Id] = element.Name
	}
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/movies", handleMovies)
	http.HandleFunc("/external", handleExternal)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
