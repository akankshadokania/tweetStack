package Handler

import (

	"net/http"
	"github.com/akankshadokania/tweetStack/model"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

	"github.com/akankshadokania/tweetStack/utils"
)


func CreateAnswer(w http.ResponseWriter, r *http.Request){

	defer r.Body.Close()
	var answer model.Answers

	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		fmt.Printf("Alert: Error")
		utils.RespondWithError(w, http.StatusBadRequest, "Error in decoding request paylod: Invalid request payload")
		return
	}

	answer.ID = bson.NewObjectId()
	if err := mongodb.Insert(answer, model.AnswerCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = answer.SetReferenceInQuestion()
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusCreated, answer)
	return

}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	obj, err := mongodb.FindById(params["id"],model.AnswerCollection)
	var answer model.Answers
	if obj == nil{
		fmt.Printf("Cannot find answer with Id %s", params["id"])
		utils.RespondWithError(w, http.StatusBadRequest, "Cannot find answer to be deleted")
		return
	}
	jsonData, err := json.Marshal(obj)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.Unmarshal(jsonData, &answer)
	err = answer.RemoveReferenceInQuestion()
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = mongodb.Delete(answer,model.AnswerCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}


func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var answer model.Answers
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Update(answer,answer.ID,model.AnswerCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindAnswer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var answer model.Answers
	obj, err := mongodb.FindById(params["id"],model.AnswerCollection)
	if obj == nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Answer with the given ID doesn't exist")
		return
	}
	jsonData, err := json.Marshal(obj)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, "Error in marshalling the answer fetched from db")
		return
	}
	err = json.Unmarshal(jsonData, &answer)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to fetch answer from DB")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, answer)
}

func AllAnswers(w http.ResponseWriter, r *http.Request) {

	objects, err := mongodb.FindAll(model.AnswerCollection)

	var answers []model.Answers
	for i := 0; i < len(objects); i++{
		var answer model.Answers
		jsonData, err := json.Marshal(objects[i])
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = json.Unmarshal(jsonData, &answer)
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		answers = append(answers, answer)
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, answers)

}

