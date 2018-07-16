package model

import "gopkg.in/mgo.v2/bson"

type Answers struct{

	ID          bson.ObjectId `bson:"_id" json:"id"`
	Answer    	string        `bson:"answer" json:"answer"`
	User        string        `bson:"user" json:"user"`
	Tags        []interface{}      `bson:"tags" json:"tags"`
	Upvotes     string        `bson:"upvotes" json:"upvotes"`
}


