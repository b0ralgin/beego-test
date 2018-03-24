package models

import (
	"errors"

	"github.com/globalsign/mgo/bson"
)

var (
	NoUser = errors.New("User Not Exist")
)

func init() {

}

type User struct {
	ID       bson.ObjectId `bson:"ID"`
	Username string        `bson:"username"`
	Password string        `bson:"password_hash"`
	Profile  *Profile      `bson:"profile"`
}

type Profile struct {
	Gender  string `bson:"gender"`
	Age     int    `bson:"age"`
	Address string `bson:"address"`
	Email   string `bson:"email"`
}

func AddUser(u *User) (string, bool) {
	u.ID = bson.NewObjectId()
	return u.ID.String(), true
}

func (u *User) UpdateProfile(profile *Profile) {
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
