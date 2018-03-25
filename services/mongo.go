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

func (m *Mongo) FindUser(uid string) (u *models.User, err error) {
	err = m.GetColl().FindId(bson.ObjectIdHex(uid)).One(&u)
	return
}

func (m *Mongo) FindUserByName(name string) (u *models.User, err error) {
	err = m.GetColl().Find(bson.M{"username": name}).One(&u)
	return
}

func (m *Mongo) UpdateProfile(user *models.User) error {
	change := bson.M{"$set": bson.M{"profile": user.Profile}}
	return m.GetColl().UpdateId(user.ID, change)
}

func (m *Mongo) UpdatePassword(user *models.User, password string) error {
	change := bson.M{"$set": bson.M{"password_hash": password}}
	return m.GetColl().UpdateId(user.ID, change)
}

func (m *Mongo) AddUser(user models.User) error {
	user.AddID()

	return m.GetColl().Insert(user)
}

func (m *Mongo) DeleteUser(uid string) error {
	return m.GetColl().RemoveId(bson.ObjectIdHex(uid))
}

func (m *Mongo) AuthenticateUser(userName string, password string) (u *models.User, err error) {
	err = m.GetColl().Find(bson.M{
		"$and": []interface{}{
			bson.M{"username": userName},
			bson.M{"password_hash": password},
		},
	}).One(&u)
	return
}

func (m *Mongo) GetColl() *mgo.Collection {
	return m.DB(beego.AppConfig.String("db")).C("Users")
}
