package Handler

import (

	"github.com/akankshadokania/tweetStack/mongo"
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
	QuestionCollection = "questions"
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
	if err := mongodb.Insert(question, QuestionCollection); err != nil {
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
	if err := mongodb.Delete(question,QuestionCollection); err != nil {
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
	if err := mongodb.Update(question,question.ID,QuestionCollection); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Print("coming here")
	var question model.Questions
	obj, err := mongodb.FindById(params["id"],QuestionCollection)
	objElem := reflect.ValueOf(obj)
	questionKeys := objElem.MapKeys()

	for _, key := range questionKeys{
		if key.String() == "_id"{
			value := objElem.MapIndex(key).Interface()
			idKey := "ID"
			err := utils.SetField(&question, idKey, value)
			if err != nil {
				fmt.Printf("Cannot set ID for struct %s\n",  err.Error())
				return
			}
			continue
		}
		fmt.Printf("The value is %s \n",objElem.MapIndex(key).Interface())
		fmt.Printf("The type of field %s is %s ,", key.String(),objElem.MapIndex(key).Kind())
		err = utils.SetField(&question, strings.Title(key.String()), objElem.MapIndex(key).Interface())
		if err != nil {
			fmt.Printf("Cannot set %s for struct %s \n", strings.Title(key.String()), err.Error())
			return
		}
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	respondWithJson(w, http.StatusOK, question)
}

func AllQuestions(w http.ResponseWriter, r *http.Request) {
	objects, err := mongodb.FindAll(QuestionCollection)

	var questions []model.Questions
	for i := 0; i < len(objects); i++{
		var question model.Questions
		objElem := reflect.ValueOf(objects[i])
		questionKeys := objElem.MapKeys()

		for _, key := range questionKeys{
			if key.String() == "_id"{
				value := objElem.MapIndex(key).Interface()
				idKey := "ID"
				err := utils.SetField(&question, idKey, value)
				if err != nil {
					fmt.Printf("Cannot set ID for struct %s",  err.Error())
					return
				}
				continue
			}
			err = utils.SetField(&question, strings.Title(key.String()), objElem.MapIndex(key).Interface())
			if err != nil {
				fmt.Printf("Cannot set %s for struct %s", strings.Title(key.String()), err.Error())
				return
			}
		}
		questions = append(questions, question)

	}


	fmt.Printf("The questions are %s", questions)
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

