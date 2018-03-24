package services

import (
	"github.com/astaxie/beego"
	"github.com/b0ralgin/test-beego/models"
	"github.com/b0ralgin/test-beego/utilities"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Mongo struct {
	*mgo.Session
}

func Startup() (*Mongo, error) {
	mongo := utilities.CopySession()
	return &Mongo{mongo}, mongo.Ping()
}

func (m *Mongo) FindUser(uid string) (*models.User, error) {
	u := models.User{}
	err := m.DB(beego.AppConfig.String("DB")).C("Users").Find(bson.M{"ID": uid}).One(&u)
	return &u, err
}

func (m *Mongo) RemoveProfile(user models.User) error {
	change := bson.M{"$unset": bson.M{"profile": nil}}
	return m.DB(beego.AppConfig.String("DB")).C("Users").Update(bson.M{"id": user.ID}, change)
}

func (m *Mongo) UpdateProfile(user *models.User, profile *models.Profile) error {
	user.UpdateProfile(profile)
	change := bson.M{"$set": bson.M{"profile": profile}}
	return m.DB(beego.AppConfig.String("DB")).C("Users").Update(bson.M{"id": user.ID}, change)
}

func (m *Mongo) AddUser(user models.User) error {
	return m.DB(beego.AppConfig.String("DB")).C("Users").Insert(user)
}

func (m *Mongo) DeleteUser(userId string) error {
	return m.DB(beego.AppConfig.String("DB")).C("Users").Remove(bson.M{"ID": userId})
}

func (m *Mongo) AuthenticateUser(userId string, password string) (bool, error) {
	n, err := m.DB(beego.AppConfig.String("DB")).C("Users").Find(bson.M{
		"$and": []interface{}{
			bson.M{"ID": userId},
			bson.M{"password_hash": password},
		},
	}).Count()
	if err != nil {
		return false, err
	}
	if n != 1 {
		return false, nil
	}
	return true, nil
}
