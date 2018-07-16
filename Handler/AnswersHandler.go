package Handler

import (

	"net/http"
	"github.com/akankshadokania/tweetStack/model"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"reflect"

	"github.com/akankshadokania/tweetStack/utils"
	"strings"
)
const(
	AnswerCollection = "answers"
)



func CreateAnswer(w http.ResponseWriter, r *http.Request){

	defer r.Body.Close()
	var answer model.Answers

	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		fmt.Printf("Alert: Error")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	fmt.Printf("Answer: %v", answer)
	answer.ID = bson.NewObjectId()
	if err := mongodb.Insert(answer, AnswerCollection); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, answer)
	return

}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var answer model.Answers
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Delete(answer,AnswerCollection); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}


func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var answer model.Answers
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Update(answer,answer.ID,AnswerCollection); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindAnswer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Print("coming here")
	var answer model.Answers
	obj, err := mongodb.FindById(params["id"],AnswerCollection)
	objElem := reflect.ValueOf(obj)
	answerKeys := objElem.MapKeys()

	for _, key := range answerKeys{
		if key.String() == "_id"{
			value := objElem.MapIndex(key).Interface()
			idKey := "ID"
			err := utils.SetField(&answer, idKey, value)
			if err != nil {
				fmt.Printf("Cannot set ID for struct %s\n",  err.Error())
				return
			}
			continue
		}

		err = utils.SetField(&answer, strings.Title(key.String()), objElem.MapIndex(key).Interface())
		if err != nil {
			fmt.Printf("Cannot set %s for struct %s \n", strings.Title(key.String()), err.Error())
			return
		}
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	respondWithJson(w, http.StatusOK, answer)
}

func AllAnswers(w http.ResponseWriter, r *http.Request) {
	objects, err := mongodb.FindAll(AnswerCollection)

	var answers []model.Answers
	for i := 0; i < len(objects); i++{
		var answer model.Answers
		objElem := reflect.ValueOf(objects[i])
		answerKeys := objElem.MapKeys()

		for _, key := range answerKeys{
			if key.String() == "_id"{
				value := objElem.MapIndex(key).Interface()
				idKey := "ID"
				err := utils.SetField(&answer, idKey, value)
				if err != nil {
					fmt.Printf("Cannot set ID for struct %s",  err.Error())
					return
				}
				continue
			}
			err = utils.SetField(&answer, strings.Title(key.String()), objElem.MapIndex(key).Interface())
			if err != nil {
				fmt.Printf("Cannot set %s for struct %s", strings.Title(key.String()), err.Error())
				return
			}
		}
		answers = append(answers, answer)

	}


	fmt.Printf("The answers are %s", answers)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, answers)

}

