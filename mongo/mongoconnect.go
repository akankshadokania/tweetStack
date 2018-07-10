package mongo

import (
	"log"
	mgdb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/akankshadokania/tweetstack/model"
	"fmt"
)

const(
	COLLECTION ="Questions"
)


type MongodbConnect struct{
	Server string
	Database string
}

var db *mgdb.Database

func (m *MongodbConnect)Connect() {

	session, err := mgdb.Dial(m.Server)
	if err != nil {
		log.Print("The error is not nil")
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}


func (m *MongodbConnect) FindById(id string) (model.Questions,error){

	var question model.Questions
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&question)

	return question, err
}

func (m *MongodbConnect) Insert(question model.Questions) error {

	fmt.Printf("Questions are %v", question)
	err := db.C(COLLECTION).Insert(&question)
	return err
}

func (m *MongodbConnect) Delete(question model.Questions) error{

	err := db.C(COLLECTION).Remove(&question)

	return err
}

func (m *MongodbConnect) Update(question model.Questions) error{

	err := db.C(COLLECTION).UpdateId(question.ID,&question)
	return err
}


func (m *MongodbConnect) FindAll() ([]model.Questions, error){

	var questions []model.Questions
	err := db.C(COLLECTION).Find(bson.M{}).All(&questions)

	return questions, err

}


