package mongo

import (
	"log"
	mgdb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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


func (m *MongodbConnect) FindById(id string, collectionName string) (interface{},error){

	var object interface{}
	err := db.C(collectionName).FindId(bson.ObjectIdHex(id)).One(&object)

	return object, err
}

func (m *MongodbConnect) Insert(obj interface{}, collectionName string) error {

	err := db.C(collectionName).Insert(&obj)
	return err
}

func (m *MongodbConnect) Delete(obj interface{}, collectionName string) error{

	err := db.C(collectionName).Remove(&obj)

	return err
}

func (m *MongodbConnect) Update(obj interface{}, ID bson.ObjectId, collectionName string) error{

	err := db.C(collectionName).UpdateId(ID,&obj)
	return err
}


func (m *MongodbConnect) FindAll(collectionName string) ([]interface{}, error){

	var objects []interface{}
	err := db.C(collectionName).Find(bson.M{}).All(&objects)

	return objects, err

}


