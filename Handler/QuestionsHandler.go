package Handler

import (

	"github.com/akankshadokania/tweetStack/mongo"
	"net/http"
	"github.com/akankshadokania/tweetStack/model"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

	"github.com/akankshadokania/tweetStack/utils"
)


var mongodb mongo.MongodbConnect

func CreateQuestion(w http.ResponseWriter, r *http.Request){

	defer r.Body.Close()
	var question model.Questions

	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		fmt.Printf("Alert: Error")
		utils.RespondWithError(w, http.StatusInternalServerError, " Error decoding request body: Invalid request payload")
		return
	}

	question.ID = bson.NewObjectId()
	if err := mongodb.Insert(question, model.QuestionCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, question)
	return

}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var question model.Questions
	obj, err := mongodb.FindById(params["id"],model.QuestionCollection)
	if obj == nil{
		utils.RespondWithError(w, http.StatusBadRequest,"Cannot find the question to be deleted")
		fmt.Printf("Cannot find question with Id %s", params["id"])
		return
	}
	jsonData, err := json.Marshal(obj)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	json.Unmarshal(jsonData, &question)
	if err := mongodb.Delete(obj,model.QuestionCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = question.RemoveRefAnswers()
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}


func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var question model.Questions
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := mongodb.Update(question,question.ID,model.QuestionCollection); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var question model.Questions
	obj, err := mongodb.FindById(params["id"],model.QuestionCollection)
	if err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if obj == nil{
		utils.RespondWithJson(w, http.StatusBadRequest, "No question found with the given ID. Question doesn't exist")
		return
	}
	jsonData, errMsg := json.Marshal(obj)
	if errMsg != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, errMsg.Error())
		return
	}
	err = json.Unmarshal(jsonData, &question)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJson(w, http.StatusOK, question)
}

func FindAllQuestions(w http.ResponseWriter, r *http.Request) {

	objects, err := mongodb.FindAll(model.QuestionCollection)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var questions []model.Questions
	for i := 0; i < len(objects); i++{
		var question model.Questions
		jsonData, err := json.Marshal(objects[i])
		if err != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		errMsg := json.Unmarshal(jsonData, &question)
		if errMsg != nil{
			utils.RespondWithError(w, http.StatusInternalServerError, errMsg.Error())
			return
		}
		questions = append(questions, question)
	}


	utils.RespondWithJson(w, http.StatusOK, questions)

}





