package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/akankshadokania/tweetStack/ModelReference"
	"k8s.io/apimachinery/pkg/util/json"
	"fmt"
	"errors"
)

const(
	AnswerCollection = "answers"
)
type Answers struct{

	ID          bson.ObjectId                  `bson:"_id" json:"_id"`
	Answer    	string        				   `bson:"answer" json:"answer"`
	User        string        				   `bson:"user" json:"user"`
	Tags        []string      			       `bson:"tags" json:"tags"`
	Upvotes     int        					   `bson:"upvotes" json:"upvotes"`
	Question    ModelReference.ModelRef        `bson:"question" json:"question"`
}


type AnswersIf interface{

	GetID () bson.ObjectId
	SetID(id bson.ObjectId)

	GetAnswer() string
	SetAnswer(answer string)

	GetUser() string
	SetUser(user string)

	GetTags() []interface{}
	SetTags(tags []interface{})

	GetUpvotes() int
	SetUpvotes(upvotes int)

	GetRefQuestion() ModelReference.ModelRef
	SetRefQuestion(question ModelReference.ModelRef)

}


func (answer Answers)GetRefQuestion() ModelReference.ModelRef{
	return answer.Question
}

func (answer *Answers) SetRefQuestion(ref ModelReference.ModelRef){
	answer.Question = ref
}

func (answer Answers)FindRefQuestions()(interface{}, error){

	questionId := answer.Question.Id
	questionCollection := answer.Question.Collection
	if len(questionId) == 0{
		fmt.Printf("Question reference is not set in Answer")
		return nil, nil
	}
	tempQuestion, err := mongodb.FindById(questionId, questionCollection)
	var question Questions
	if err != nil{
		return nil, err
	}
	jsonData, err := json.Marshal(tempQuestion)
	if err != nil{
		return nil,err
	}
	errMsg := json.Unmarshal(jsonData, &question)
	if errMsg != nil{
		return nil, errMsg
	}
	 return question, nil
}

func(answer Answers) SetReferenceInQuestion() error{

	answerRef := ModelReference.ModelRef{Id : answer.ID.Hex(), Collection: AnswerCollection, Db: ""}
	questionInt, err := answer.FindRefQuestions()
	if err != nil {
		return err
	}
	question := questionInt.(Questions)
	oldAnsRefs := question.GetRefAnswers()
	oldAnsRefs= append(oldAnsRefs, answerRef)
	question.Answers = oldAnsRefs
	updateError := mongodb.Update(question,question.ID, answer.Question.Collection)
	if updateError != nil{
		return updateError
	}
	return nil
}


func (answer Answers) RemoveReferenceInQuestion()error{

	questionInt, err := answer.FindRefQuestions()
	if err != nil{
		return err
	}

	if questionInt == nil{
		errMsg := errors.New("Cannot find any referenced question in Answer")
		return errMsg
	}

	question := questionInt.(Questions)
	refAnswers := question.GetRefAnswers()
	if len(refAnswers) != 0{
		for i:=0; i< len(refAnswers); i++{
			if refAnswers[i].Id == answer.ID.Hex(){
				refAnswers[i]=refAnswers[len(refAnswers)-1]
				refAnswers = refAnswers[:len(refAnswers)-1]
				break
			}
		}
	}
	question.Answers=refAnswers
	updateError := mongodb.Update(question,question.ID, answer.Question.Collection)
	if updateError != nil{
		return updateError
	}
	return nil

}


