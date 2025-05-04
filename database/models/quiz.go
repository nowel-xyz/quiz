package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Quiz struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Questions []Question         `bson:"questions"`
}

type Question struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Text     string             `bson:"text"`
	Options  []string           `bson:"options"`
	Answer   string             `bson:"answer"`
}
  