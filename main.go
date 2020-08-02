package main

import (
	"encoding/json" //for jason marshalling and unmarshalling
	"fmt"
	"io/ioutil" //for reading response body
	"log"
	"net/http"
	"strconv"

	mssqlserver "my-rest-project/sqlservconnect" //For Database CRUD Operations

	"github.com/gorilla/mux" // router
)

//defining resource type
type Movie struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

func main() {
	
	mssqlserver.StartDatabaseServer() //start SQL SERVER

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

	// Read Movies
	Movies, err := mssqlserver.ReadMovies()
	if err != nil {
		log.Fatal("Error reading Movies: ", err.Error())
	}
	json.NewEncoder(w).Encode(Movies) //encode response in json and write to response Writer
}

func getSingleMovie(w http.ResponseWriter, r *http.Request) {
	//Get specific resource Handler function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters

	movie, err := mssqlserver.ReadSingleMovie(id) //Read Specific record from database
	if err != nil{
		fmt.Fprintf(w, "Movie not found")
		log.Fatal("Error reading Movies: ", err.Error())
	}
	if movie != nil {
		json.NewEncoder(w).Encode(movie) //encode response in json and write to response Writer
	} else {
		fmt.Fprintf(w, "Movie not found")
	}
}

func postMovie(w http.ResponseWriter, r *http.Request) {
	//Post resource Handler function
	var movie Movie
	reqBody, _ := ioutil.ReadAll(r.Body) // get the body of our POST request
	json.Unmarshal(reqBody, &movie)      // unmarshal this into a new movie struct

	err := mssqlserver.CreateMovie(movie.Name, movie.Year)  //insert movie into database
	if err != nil {
		fmt.Fprintf(w, "Error")
		log.Fatal("Error while inserting movie: ", err.Error())
	}
	fmt.Fprintf(w, "Movie Inserted Succesfully!")
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	//Delete resource Handler Function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters

	err := mssqlserver.DeleteMovie(id) // Delete movie from database
	if err != nil {
		fmt.Fprintf(w, "Error")
		log.Fatal("Error while deleting movie: ", err.Error())
	}
	fmt.Fprintf(w, "Movie Deleted Succesfully!")
}

func putMovie(w http.ResponseWriter, r *http.Request) {
	//Update resource Handler Function
	variables := mux.Vars(r)               //parse the path parameters
	id, _ := strconv.Atoi(variables["id"]) //extract id from path parameters
	reqBody, _ := ioutil.ReadAll(r.Body)   //read request body

	var newMovie Movie
	json.Unmarshal(reqBody, &newMovie) //unmarshal request into movie struct

	err := mssqlserver.UpdateMovie(id,newMovie.Name, newMovie.Year)  //update movie in database
	if err != nil {
		fmt.Fprintf(w, "Error")
		log.Fatal("Error while updating movie: ", err.Error())
	}
	fmt.Fprintf(w, "Movie Updated Succesfully!")
}
