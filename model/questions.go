package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/akankshadokania/tweetStack/ModelReference"
	"github.com/akankshadokania/tweetStack/mongo"
	"k8s.io/apimachinery/pkg/util/json"
	"fmt"
)

var mongodb mongo.MongodbConnect

const(
	QuestionCollection = "questions"
)

type Questions struct{

	ID          bson.ObjectId 				`bson:"_id" json:"_id"`
	Question    string        				`bson:"question" json:"question"`
	User        string        				`bson:"user" json:"user"`
    Upvotes     int       	 				`bson:"upvotes" json:"upvotes"`
	Tags        []string     				`bson:"tags" json:"tags"`

	Answers     []ModelReference.ModelRef   `bson:"answers" json:"answers"`
}


type QuestionsIf interface{

	GetID()bson.ObjectId
	SetID(id bson.ObjectId)

	GetQuestion()string
	SetQuestion(question string)

	GetUser()string
	SetUser(user string)

	GetUpvotes()int
	SetUpvotes(upvotes int)

	GetTags()[]interface{}
	SetTags([]interface{})

	GetRefAnswers()ModelReference.ModelRef
	SetRefAnswers(ref ModelReference.ModelRef)

}


func (question Questions)FindRefAnswers()([]Answers, error){

	ansRef := question.GetRefAnswers()

	if len(ansRef) == 0 {
		return nil, nil
	}

	var answers []Answers
	var answer  Answers

	for _, ans := range ansRef {
		tempAnswer, err := mongodb.FindById(ans.Id, ans.Collection)
		if err != nil {
			return nil, err
		}
		jsonData, err := json.Marshal(tempAnswer)
		if err == nil {
			json.Unmarshal(jsonData, &answer)
			answers = append(answers, answer)
		} else {
			return nil, err
		}
	}

	return answers, nil

}


func (question Questions)GetRefAnswers() []ModelReference.ModelRef{
	return question.Answers
}

func (question Questions)GetRefCollectionName() string{
	refAns := question.GetRefAnswers()
	if len(refAns) > 0{
		return refAns[0].Collection
	}
	return ""
}

func (question Questions)RemoveRefAnswers() error{

	refAns, err := question.FindRefAnswers()
	if err != nil {
		fmt.Printf("Cannot fetch remote answers. Error %s", err.Error())
		return err
	}
	if len(refAns) == 0{
		return nil
	}
	collectionName := question.GetRefCollectionName()
	for _, tempAns := range refAns{
		err := mongodb.Delete(tempAns, collectionName)
		if err != nil {
			fmt.Printf("Cannot delete answer %s for question %s ,Error %s", tempAns.Answer, question.Question, err.Error() )
			return err
		}
	}
	return nil
}

