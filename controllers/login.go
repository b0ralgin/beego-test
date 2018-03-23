package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/b0ralgin/test-beego/models"
	jwt "github.com/dgrijalva/jwt-go"
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
func (o *LoginController) SignIn() {
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
// @Failure 409 user is exist
// @router /v1/login [post]
func (o *LoginController) SignUp() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	if id, ok := models.AddUser(user); ok {
		o.CustomAbort(409, "User already exist")
	} else {
		user.Id = id
	}
	o.Data["json"] = user
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
