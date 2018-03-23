package controllers

import (
	"encoding/json"

	"github.com/b0ralgin/test-beego/models"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/astaxie/beego"
)

// Operations about object
type LoginController struct {
	beego.Controller
}

const SIGNKEY = "B1rljsRlQlH7+NSvWuFjU/DROpULnTFB"

type customClaims struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

// @Title Login
// @Description provides  authentication of user
// @Param	body	body 	models.User	true		"The object content"
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 403 user not exist
// @Failure 500 errors in function
// @router /v1/login [post]
func (o *LoginController) Login() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	if !models.Login(user.Id, user.Password) {
		o.CustomAbort(403, "Wrong username of password")
	} else {
		o.Data["json"], err = generateJWTToken(user)
		if err != nil {
			o.CustomAbort(500, err.Error())
		}
	}
	o.ServeJSON()
}

// @Title Sign
// @Description provides  creation of user
// @Param	body	body 	models.User	true		"The object content"
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 403 user not exist
// @Failure 500 errors in function
// @router /v1/login [post]
func (o *LoginController) Sign() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	if user, err := models.GetUser(user.Id); err != nil {
		if err == models.NoUser {
			o.CustomAbort(403, err.Error())
		} else {
			o.CustomAbort(500, err.Error())
		}
	} else {
		o.Data["json"], err = generateJWTToken(*user)
		if err != nil {
			o.CustomAbort(500, err.Error())
		}
	}
	o.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *LoginController) Get() {
	objectId := o.Ctx.Input.Param(":objectId")
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *LoginController) GetAll() {
	obs := models.GetAll()
	o.Data["json"] = obs
	o.ServeJSON()
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *LoginController) Put() {
	objectId := o.Ctx.Input.Param(":objectId")
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		o.Data["json"] = err.Error()
	} else {
		o.Data["json"] = "update success!"
	}
	o.ServeJSON()
}

// @Title Delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (o *LoginController) Delete() {
	objectId := o.Ctx.Input.Param(":objectId")
	models.Delete(objectId)
	o.Data["json"] = "delete success!"
	o.ServeJSON()
}

func generateJWTToken(user models.User) (map[string]string, error) {
	claims := &customClaims{
		ID: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SIGNKEY))
	return map[string]string{"token": tokenString}, err
}
