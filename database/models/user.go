package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPEntry struct {
	IP         string    `bson:"ip"`
	LoginTimes int       `bson:"loginTimes"`
	LastLogin  time.Time `bson:"lastLogin"`
}

type User struct {
	ObjID   primitive.ObjectID `bson:"_id,omitempty"`
	ID      string             `bson:"id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Salt     string             `bson:"salt"`
	Email    string             `bson:"email"`
	Cookie   string             `bson:"cookie"`
	IPs      []IPEntry          `bson:"ips"`
	Roles    []string           `bson:"roles"`
}


type LobbyUser struct {
	ID   string             `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Roles    []string           `json:"roles"`
}