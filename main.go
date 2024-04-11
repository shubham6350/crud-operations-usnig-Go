package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct{
	Firstname string `json:"firstname"`
	Lastname string `json: "lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for index,item := range movies{
		if item.ID==params["id"]{
			movies=append( movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovieByID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params:= mux.Vars(r)
	for _,item:=range movies{
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return 
		}
	}
	//If the movie with given ID is not found send an HTTP response status code 404 and a message saying that the movie was
	err:=errors.New("Movie not found.")
	http.Error(w, err.Error(), 404)
}

func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID=strconv.Itoa(rand.Intn(10000000))
	movies=append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovieByID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params:=mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies=append(movies[:index],movies[index+1:]...)
			var movie Movie
			_=json.NewDecoder(r.Body).Decode(&movie)
			movie.ID=params["id"]
			movies=append(movies,movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	err:= errors.New("Movie not found.")
	http.Error(w, err.Error(), http.StatusNotFound)
}


func main(){
	r:=mux.NewRouter()
	movies = append(movies,Movie{"1","9780", "The Lord of the Rings - The Return of the King", &Director {"Peter", "Jackson"}})
	movies = append(movies,Movie{"2","9781","wanted",&Director{"karan","johar"}})
	movies = append(movies,Movie{"3","8895","Family man",&Director{"Raj","DK"}})
	movies = append(movies,Movie{"4","4455","Farzi",&Director{"sahid","kapoor"}})

	r.HandleFunc("/movies",getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}",getMovieByID).Methods("GET")
	r.HandleFunc("/movies",createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}",updateMovieByID).Methods("PUT")
	r.HandleFunc("/movies/{id}",deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at Port 8000\n")
	log.Fatal(http.ListenAndServe( ":8000", r))
}