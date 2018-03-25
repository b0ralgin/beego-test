package models

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
)

func init() {

}

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"-"`
	Username string        `bson:"username"`
	Password passwordHash  `bson:"password_hash"`
	Profile  *Profile      `bson:"profile"`
}

type Profile struct {
	Gender  string `bson:"gender"`
	Age     int    `bson:"age"`
	Address string `bson:"address"`
	Email   string `bson:"email"`
}

type passwordHash string

func (u *User) AddID() {
	u.ID = bson.NewObjectId()
}

func (u *User) UpdateProfile(profile Profile) {
	if u.Profile == nil {
		u.Profile = new(Profile)
	}
	if profile.Age != 0 {
		u.Profile.Age = profile.Age
	}
	if profile.Address != "" {
		u.Profile.Address = profile.Address
	}
	if profile.Gender != "" {
		u.Profile.Gender = profile.Gender
	}
	if profile.Email != "" {
		u.Profile.Email = profile.Email
	}
}

func (p passwordHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
