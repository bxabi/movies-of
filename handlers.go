package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type Page struct {
	Data     interface{}
	NextPage string
}

type Actor struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Index", &Page{Data: "Movies of"})
}

type SearchResponse struct {
	Page       int     `json:"page"`
	TotalPages int     `json:"total_pages"`
	Actors     []Actor `json:"results"`
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	u, _ := url.Parse("https://api.themoviedb.org/3/search/person?include_adult=false&language=en-US")
	q := u.Query()
	q.Add("api_key", apiKey)
	q.Add("query", term)
	q.Add("page", page)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode != 200 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var result SearchResponse
	json.NewDecoder(resp.Body).Decode(&result)

	var nextPage string
	if result.Page > 0 && result.Page < result.TotalPages {
		next := url.Values{}
		next.Set("term", term)
		next.Set("page", strconv.Itoa(result.Page+1))
		nextPage = "/search?" + next.Encode()
	}

	renderTemplate(w, "ActorList", &Page{
		Data:     result.Actors,
		NextPage: nextPage,
	})
}

type Movie struct {
	Title        string
	Release_date string
	Vote_count   int
	Vote_average float32
	Poster_path  string
	Character    string
	Overview     string
	Media_type   string
	Name         string
	Credit_id    string
	Genre_ids    []int
	Id           int
}

type MoviesPage struct {
	Cast   []Movie
	Genres map[int]string
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	actor := r.URL.Query().Get("actor")
	u, _ := url.Parse("https://api.themoviedb.org/3/person/" + actor + "/combined_credits?language=en-US")
	q := u.Query()
	q.Add("api_key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var result map[string][]Movie
	json.NewDecoder(resp.Body).Decode(&result)

	array := result["cast"]

	sort.Slice(array, func(i, j int) bool {
		return array[i].Vote_average > array[j].Vote_average
	})

	renderTemplate(w, "Movies", &Page{
		Data: MoviesPage{Cast: array, Genres: Genres},
	})
}

func handleExternal(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	typ := r.URL.Query().Get("type")
	u, _ := url.Parse("https://api.themoviedb.org/3/" + typ + "/" + id + "/external_ids?")
	q := u.Query()
	q.Add("api_key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		// c.Flash.Error("Error during retrieving the link: " + resp.Status)
		// c.FlashParams()
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	http.Redirect(w, r, "https://www.imdb.com/title/"+result["imdb_id"], http.StatusFound)
}
