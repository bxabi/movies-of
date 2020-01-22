package controllers

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

type Genre struct {
	Id   int
	Name string
}

var apiKey string

var Genres = make(map[int]string)

func init() {
	loadApiKey()
	loadGenres()
}

func loadApiKey() {
	pwd, _ := os.Getwd()
	// this points to the app folder.
	file, err := os.Open(pwd + "/apiKey")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	content := scanner.Text()

	apiKey = content
}

func loadGenres() {
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

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Search(term string) revel.Result {
	u, _ := url.Parse("https://api.themoviedb.org/3/search/person?page=1&include_adult=false&language=en-US")
	q := u.Query()
	q.Add("api_key", apiKey)
	q.Add("query", term)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode != 200 {
		c.Flash.Error("Error during the search: " + resp.Status)
		c.FlashParams()
		return c.Redirect("/")
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.ViewArgs["results"] = result["results"]
	c.ViewArgs["term"] = term
	return c.RenderTemplate("App/ActorList.html")
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

func (c App) Movies(actor string) revel.Result {
	u, _ := url.Parse("https://api.themoviedb.org/3/person/" + actor + "/combined_credits?language=en-US")
	q := u.Query()
	q.Add("api_key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		c.Flash.Error("Error retrieving the movies: " + resp.Status)
		c.FlashParams()
		return c.Redirect("/")
	}
	defer resp.Body.Close()

	var result map[string][]Movie
	json.NewDecoder(resp.Body).Decode(&result)

	array := result["cast"]

	sort.Slice(array, func(i, j int) bool {
		return array[i].Vote_average > array[j].Vote_average
	})

	c.ViewArgs["cast"] = array
	c.ViewArgs["genres"] = Genres
	return c.RenderTemplate("App/Movies.html")
}

func (c App) External(id string) revel.Result {
	u, _ := url.Parse("https://api.themoviedb.org/3/movie/" + id + "/external_ids?")
	q := u.Query()
	q.Add("api_key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		c.Flash.Error("Error during retrieving the link: " + resp.Status)
		c.FlashParams()
		return c.Redirect("/")
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	c.ViewArgs["ids"] = result
	return c.Redirect("https://www.imdb.com/title/" + result["imdb_id"])
}
