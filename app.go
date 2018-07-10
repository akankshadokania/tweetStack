package main

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"log"
	"github.com/akankshadokania/tweetstack/model"
	"github.com/akankshadokania/tweetstack/mongo"
	"fmt"

)

var mongodb mongo.MongodbConnect

func CreateQuestion(w http.ResponseWriter, r *http.Request){

	defer r.Body.Close()
	var question model.Questions

	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		fmt.Printf("Alert: Error")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	fmt.Printf("Questions: %v", question)
	question.ID = bson.NewObjectId()
	if err := mongodb.Insert(question); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, question)
	return

}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var question model.Questions
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Delete(question); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}


func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var question model.Questions
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Update(question); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := mongodb.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func AllQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := mongodb.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, questions)
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {

	mongodb.Database= "mongo_cont"
	mongodb.Server = "localhost"
	mongodb.Connect()

}

// Define HTTP request routes
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/Questions", AllQuestions).Methods("GET")
	r.HandleFunc("/Questions", CreateQuestion).Methods("POST")
	r.HandleFunc("/Questions", UpdateQuestion).Methods("PUT")
	r.HandleFunc("/Questions", DeleteQuestion).Methods("DELETE")
	r.HandleFunc("/Questions/{id}", FindQuestion).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
