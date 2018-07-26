package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/akankshadokania/tweetStack/mongo"
	"github.com/akankshadokania/tweetStack/Handler"
)

var mongodb mongo.MongodbConnect


func init() {

	mongodb.Database= "mongo_cont"
	mongodb.Server = "localhost"
	mongodb.Connect()

}

// Define HTTP request routes
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/Questions", Handler.FindAllQuestions).Methods("GET")
	r.HandleFunc("/Questions", Handler.CreateQuestion).Methods("POST")
	r.HandleFunc("/Questions", Handler.UpdateQuestion).Methods("PUT")
	r.HandleFunc("/Questions/{id}", Handler.DeleteQuestion).Methods("DELETE")
	r.HandleFunc("/Questions/{id}", Handler.FindQuestion).Methods("GET")
	r.HandleFunc("/Answers", Handler.AllAnswers).Methods("GET")
	r.HandleFunc("/Answers", Handler.CreateAnswer).Methods("POST")
	r.HandleFunc("/Answers", Handler.UpdateAnswer).Methods("PUT")
	r.HandleFunc("/Answers/{id}", Handler.DeleteAnswer).Methods("DELETE")
	r.HandleFunc("/Answers/{id}", Handler.FindAnswer).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
