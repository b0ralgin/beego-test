package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/b0ralgin/test-beego/models"
	"github.com/b0ralgin/test-beego/services"
	jwt "github.com/dgrijalva/jwt-go"
)

// Operations about object
type LoginController struct {
	beego.Controller
	*services.Mongo
}

type customClaims struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

func (o *LoginController) Prepare() {
	var err error
	o.Mongo, err = services.Startup()
	if err != nil {
		o.CustomAbort(500, err.Error())
	}
}

// @Title SignIn
// @Description provides  authentication of user
// @Param	body	body 	models.User	true		"The object content"
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 403 user not exist
// @Failure 500 errors in function
// @router /v1/signin [post]
func (o *LoginController) SignIn() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	if ok, err := o.Mongo.AuthenticateUser(user.ID.String(), user.Password); !ok {
		if err != nil {
			o.CustomAbort(500, err.Error())
		} else {
			o.CustomAbort(403, "Wrong username of password")
		}

	} else {
		o.Data["json"], err = generateJWTToken(user)
		if err != nil {
			o.CustomAbort(500, err.Error())
		}
	}
	o.ServeJSON()
}

// @Title SignUp
// @Description provides  creation of user
// @Param	body	body 	models.User	true		"The object content"
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 409 user is exist
// @router /v1/signup [post]
func (o *LoginController) SignUp() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	if err := o.Mongo.AddUser(user); err != nil {
		o.CustomAbort(400, "Can't Parse body")
	}
	o.Data["json"] = user
	o.ServeJSON()
}

func generateJWTToken(user models.User) (map[string]string, error) {
	claims := &customClaims{
		ID: user.ID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("JWTSignKey")))
	return map[string]string{"token": tokenString}, err
}
