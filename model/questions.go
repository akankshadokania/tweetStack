package model

import "gopkg.in/mgo.v2/bson"

type Questions struct{

	ID          bson.ObjectId `bson:"_id" json:"id"`
	Question    string        `bson:"question" json:"question"`
	User        string        `bson:"user" json:"user"`
	AnswerId    string        `bson:"answerId" json:"answerId"`
    Upvotes     string        `bson:"upvotes" json:"upvotes"`
	Tags        []interface{}     `bson:"tags" json:"tags"`
}

