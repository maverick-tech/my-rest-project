package main

import (
	"encoding/json" //for jason marshalling and unmarshalling
	"fmt"
	"io/ioutil" //for reading response body
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // router
)

//defining resource type
type Movie struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

var Movies []Movie

func main() {
	//sample database
	Movies = []Movie{
		Movie{Id: 1, Name: "Avengers", Year: 2016},
		Movie{Id: 2, Name: "Avengers:Endgame", Year: 2019},
		Movie{Id: 3, Name: "Doctor Strange", Year: 2017},
		Movie{Id: 4, Name: "Ironman", Year: 2010},
		Movie{Id: 5, Name: "Thor Ragnarok", Year: 2018},
	}

	route() //handle requests
}

func route() {
	router := mux.NewRouter().StrictSlash(true)

	//Handle Read Operations
	router.HandleFunc("/", mainPage).Methods("GET")
	router.HandleFunc("/movies", getAllMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getSingleMovie).Methods("GET")

	//Handle Create Operation
	router.HandleFunc("/movies", postMovie).Methods("POST")

	//Handle Delete Operation
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	//Handle Update Operation
	router.HandleFunc("/movies/{id}", putMovie).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router)) //start to listen and serve
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my webapp!") //Main page Handler function
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	//Get resource Handler function
	json.NewEncoder(w).Encode(Movies) //encode response in json and write to response Writer
}

func getSingleMovie(w http.ResponseWriter, r *http.Request) {
	//Get specific resource Handler function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters

	for _, movie := range Movies { //loop through all entries
		if movie.Id == id { //find specific movie
			json.NewEncoder(w).Encode(movie)
			return //return that specific movie
		}
	}
	fmt.Fprintf(w, "Movie not found")
}

func postMovie(w http.ResponseWriter, r *http.Request) {
	//Post resource Handler function
	var movie Movie
	reqBody, _ := ioutil.ReadAll(r.Body) // get the body of our POST request
	json.Unmarshal(reqBody, &movie)      // unmarshal this into a new movie struct

	Movies = append(Movies, movie) //append to current database
	json.NewEncoder(w).Encode(movie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	//Delete resource Handler Function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters

	//loop through all our movies
	for index, movie := range Movies {
		// if our id path parameter matches one of our movies
		if movie.Id == id {
			//remove the movie
			Movies = append(Movies[:index], Movies[index+1:]...)
			fmt.Fprintf(w, "Movie Deleted")
			return
		}
	}
	fmt.Fprintf(w, "Movie not found !")
}

func putMovie(w http.ResponseWriter, r *http.Request) {
	//Update resource Handler Function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters
	reqBody, _ := ioutil.ReadAll(r.Body)   //read request body

	var newMovie Movie
	json.Unmarshal(reqBody, &newMovie) //unmarshal request into movie struct

	//loop through all our articles
	for index, movie := range Movies {
		if movie.Id == id { //update the movie if id found
			Movies[index].Id = newMovie.Id
			Movies[index].Name = newMovie.Name
			Movies[index].Year = newMovie.Year
			fmt.Fprintf(w, "Movie Updated\n")
			json.NewEncoder(w).Encode(Movies[index]) //encode response in json and return to response Writer
			return
		}
	}
	fmt.Fprintf(w, "Movie not found !")
}
