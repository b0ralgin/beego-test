package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/b0ralgin/test-beego/models"
	"github.com/b0ralgin/test-beego/services"
	"github.com/chilts/sid"
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

var reqCache cache.Cache

func init() {
	reqCache, _ = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":120}`)
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
	u, err := o.Mongo.AuthenticateUser(user.Username, string(user.Password))
	if u == nil {
		o.CustomAbort(403, "Wrong username of password")
	}
	if err != nil {
		o.CustomAbort(500, err.Error())
	}
	o.Data["json"], err = generateJWTToken(*u)
	if err != nil {
		o.CustomAbort(500, err.Error())
	}
	o.ServeJSON()
}

// @Title SignUp
// @Description provides  creation of user
// @Param	body	body 	models.User	true		"The object content"
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 500 failure with mongo
// @router /v1/signup [post]
func (o *LoginController) SignUp() {
	var user models.User
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &user)
	if err != nil {
		o.CustomAbort(400, "Can't Parse body: "+err.Error())
	}
	user.AddID()
	if err := o.Mongo.AddUser(user); err != nil {
		o.CustomAbort(500, err.Error())
	}
	o.Data["json"] = user
	o.ServeJSON()
}

// @Title Reset Password
// @Description provides  creation of user
// @Success 200 {string} JWTToken
// @Failure 400 body is wrong
// @Failure 409 user is exist
// @router /reset  [post]
func (o *LoginController) Reset() {
	name := o.Input().Get("name")
	requestid := o.Input().Get("requestid")
	fmt.Println(name)
	u, err := o.Mongo.FindUserByName(name)
	if u == nil {
		o.CustomAbort(400, "Wrong user")
	}
	if err != nil {
		o.CustomAbort(500, err.Error())
	}

	if requestid == "" {
		o.generateRequest(u)
	} else {
		o.makeNewPassword(u, requestid)
	}
	o.ServeJSON()
}

func (o *LoginController) Finish() {
	defer func() {
		o.Mongo.Close()
	}()
}

func generateJWTToken(user models.User) (map[string]string, error) {
	claims := &customClaims{
		ID: user.ID.Hex(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("JWTSignKey")))
	return map[string]string{"token": tokenString}, err
}

func (o *LoginController) makeNewPassword(user *models.User, requestid string) {
	if reqCache.Get(user.ID.Hex()) == nil {
		o.CustomAbort(403, "Wrong request")
	}
	if !(reqCache.Get(user.ID.Hex()).(string) == requestid) {
		o.CustomAbort(403, "Wrong request")
	}
	var password struct {
		Password string
	}
	json.Unmarshal(o.Ctx.Input.RequestBody, &password)

	err := o.Mongo.UpdatePassword(user, password.Password)
	if err != nil {
		o.CustomAbort(500, err.Error())
	}
	if err := reqCache.Delete(user.ID.Hex()); err != nil {
		o.CustomAbort(500, err.Error())
	}
	o.Data["json"] = user
}

func (o *LoginController) generateRequest(user *models.User) {
	requestid := sid.Id()
	err := reqCache.Put(user.ID.Hex(), requestid, 0)

	if err != nil {
		o.CustomAbort(500, err.Error())
	}
	o.Data["json"] = map[string]string{
		"request": requestid,
	}
}
