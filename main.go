package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

var Movies []Movie

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my webapp!")
}

func main() {
	//sample database
	Movies = []Movie{
		Movie{Id: 1, Name: "Avengers", Year: 2016},
		Movie{Id: 2, Name: "Avengers:Endgame", Year: 2019},
		Movie{Id: 3, Name: "Doctor Strange", Year: 2017},
		Movie{Id: 4, Name: "Ironman", Year: 2010},
		Movie{Id: 5, Name: "Thor Ragnarok", Year: 2018},
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", mainPage)
	router.HandleFunc("/movies", getAllMovies)
	router.HandleFunc("/movies/{id}", getSingleMovie)
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Movies)
}

func getSingleMovie(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	id, _ := strconv.Atoi(variables["id"])

	for _, movie := range Movies {
		if movie.Id == id {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	fmt.Fprintf(w, "Movie not found")
}
