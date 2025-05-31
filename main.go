package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	ISBN     string    `json:"isbn"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)

	for index, movie := range movies {
		if movie.ID == param["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)

}

func createMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	//params
	params := mux.Vars(r)
	//loop over the movies
	for index, item := range movies {
		if item.ID == params["id"] {
			// delete the movie with id received
			movies = append(movies[:index], movies[index+1:]...)
			// add a new movie - the movie we sent in the body
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			return

		}
	}

}

func initializeMovie() {
	movies = append(movies, Movie{ID: "1", Title: "Movie One", ISBN: "438748374", Director: &Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", Title: "Movie Two", ISBN: "438748375", Director: &Director{FirstName: "Jane", LastName: "Doe"}})
}

func main() {
	r := mux.NewRouter()
	initializeMovie()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Println("Starting server on port 8080....")
	log.Fatal(http.ListenAndServe(":8080", r))
}
