package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	CONN_HOST      = "localhost"
	CONN_PORT      = "8080"
	ADMIN_USER     = "admin"
	ADMIN_PASSWORD = "admin"
	API_KEY        = "?api_key=6f978b42bc40e53d6600747648d6b4a1"
	API_ROOT       = "https://api.themoviedb.org/3/"
)

type greet struct {
	Message string `json:"message"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GenresResult struct {
	Genres []Genre `json:"genres"`
}

type MovieResult struct {
	PosterPath       string  `json:"poster_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	GenreId          []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Title            string  `json:"title"`
	BackdropPath     string  `json:"backdrop_path"`
	Popularity       float32 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
}

var mapGenres map[string]string

type MovieResponse struct {
	Page         int           `json:"page"`
	MovieResults []MovieResult `json:"results"`
	TotalResults int           `json:"total_results"`
	TotalPages   int           `json:"total_pages"`
}

type MovieList struct {
	Titulo       string
	Id           int
	Genero       []string
	FechaEstreno string
	Sinopsis     string
}

type ErrResponse struct {
	Error  string
	Result int
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	var greeting greet
	greeting.Message = "Hello World!!!!"
	json.NewEncoder(w).Encode(greeting)

}

func getGenres(lang string) int {
	queryUrl := API_ROOT + "genre/movie/list" + API_KEY + "&language=" + lang
	response, err := http.Get(queryUrl)
	if err != nil {
		//log.Fatal(err)
		return -1
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.Fatal(err)
		return -2
	}

	var responseObject GenresResult
	json.Unmarshal(responseData, &responseObject)

	mapGenres = make(map[string]string)

	for _, v := range responseObject.Genres {
		id := strconv.Itoa(v.Id)
		mapGenres[id] = v.Name
	}
	return 0
}

func getMovies(MovieName string, MovieLanguage string) (ErrResponse, []MovieList) {
	var lang string
	var defaultResponse []MovieList

	switch MovieLanguage {
	case "en":
		lang = "en-US"
	case "es":
		lang = "es-ES"
	}
	var errResult ErrResponse
	if x := getGenres(lang); x != 0 {

		switch x {

		case -1:
			errResult.Result = 1
			errResult.Error = "Ocurrio un error al obtener los generos de las peliculas"
		case -2:
			errResult.Result = 2
			errResult.Error = "Ocurrio un error al  procesar los generos de las peliculas"
		}

		return errResult, defaultResponse
	}
	//eventID := mux.Vars(r)["id"]
	//fmt.Fprintf(w, "Category: %v\n", vars["category"])
	queryMovie := "&query=" + MovieName + "&page=1&include_adult=false"
	queryLanguage := "&language=" + MovieLanguage
	queryURL := API_ROOT + "search/movie" + API_KEY + queryLanguage + queryMovie

	response, err := http.Get(queryURL)
	if err != nil {
		//log.Fatal(err)
		errResult.Result = 3
		errResult.Error = "Ocurrio un error al obtener las peliculas"
		return errResult, defaultResponse
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//log.Fatal(err)
		errResult.Result = 3
		errResult.Error = "Ocurrio un error al procesar las peliculas"
		return errResult, defaultResponse

	}

	var responseObject MovieResponse
	json.Unmarshal(responseData, &responseObject)

	movieList := make([]MovieList, 1)

	for _, v := range responseObject.MovieResults {
		var m MovieList

		m.FechaEstreno = v.ReleaseDate
		for _, val := range v.GenreId {
			genreId := strconv.Itoa(val)
			m.Genero = append(m.Genero, mapGenres[genreId])
		}
		m.Id = v.Id
		m.Sinopsis = v.Overview
		m.Titulo = v.Title
		//id := strconv.Itoa(v.GenreId)
		movieList = append(movieList, m)

	}
	return errResult, movieList

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieName := vars["movieName"]
	movieLanguage := vars["language"]

	err, lista := getMovies(movieName, movieLanguage)
	if err.Result != 0 {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(lista)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HelloWorldHandler)
	router.HandleFunc("/peliculas/{movieName}/{language}", ApiHandler)
	http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
}
