package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/b0ralgin/test-beego/models"
	"github.com/b0ralgin/test-beego/services"
	jwt "github.com/dgrijalva/jwt-go"
)

// Operations about Users
type UserController struct {
	beego.Controller
	*services.Mongo
}

func (u *UserController) Prepare() {
	authorizationHeader := u.Ctx.Request.Header.Get("Authorization")
	if authorizationHeader == "" {
		u.CustomAbort(401, "Nonauthorized")
	}
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) == 2 {
		var userClaims customClaims
		token, err := jwt.ParseWithClaims(bearerToken[1], &userClaims, parseToken)
		if err != nil {
			u.CustomAbort(401, "Can't parse JWT token")
		}
		if !token.Valid {
			u.CustomAbort(401, "Nonauthorized")
		}
		u.Ctx.Input.SetData("userId", token.Claims.(*customClaims).ID)
		u.Mongo, err = services.Startup()
		if err != nil {
			u.CustomAbort(500, err.Error())
		}
	}
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	userId, ok := u.Ctx.Input.GetData("userId").(string)
	if !ok {
		u.CustomAbort(400, "wrong ID")
	}
	uid := u.GetString(":uid")
	if uid != "" {
		userId = uid
	}
	user, err := u.Mongo.FindUser(userId)
	if err != nil {
		if err != nil {
			if err == models.NoUser {
				u.CustomAbort(400, "User doesn't exist")
			} else {
				u.CustomAbort(500, err.Error())
			}
		}
		u.CustomAbort(500, err.Error())
	}

	u.Data["json"] = user
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router / [put]
func (u *UserController) Put() {
	uid, ok := u.Ctx.Input.GetData("userId").(string)
	if !ok {
		u.CustomAbort(400, "wrong ID")
	}
	user, err := u.Mongo.FindUser(uid)
	if err != nil {
		if err == models.NoUser {
			u.CustomAbort(400, "User doesn't exist")
		} else {
			u.CustomAbort(500, err.Error())
		}
	}
	var newUser models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &newUser)
	if newUser.Password != "" {
		err = u.Mongo.UpdatePassword(user, string(newUser.Password))
		if err != nil {
			u.CustomAbort(500, err.Error())
		}
	}
	if newUser.Profile != nil {
		user.UpdateProfile(*newUser.Profile)
		err = u.Mongo.UpdateProfile(user)
		if err != nil {
			u.CustomAbort(500, err.Error())
		}
	}
	u.Data["json"] = user
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router / [delete]
func (u *UserController) Delete() {
	uid, ok := u.Ctx.Input.GetData("userId").(string)
	if !ok {
		u.CustomAbort(400, "wrong ID")
	}
	if uid != "" {
		err := u.Mongo.DeleteUser(uid)
		if err != nil {
			if err == models.NoUser {
				u.CustomAbort(400, "User doesn't exist")
			} else {
				u.CustomAbort(500, err.Error())
			}
		}
		u.Data["json"] = "delete success!"
	}
	u.ServeJSON()
}

func parseToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}
	return []byte(beego.AppConfig.String("JWTSignKey")), nil
}
