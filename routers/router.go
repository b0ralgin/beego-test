// @APIVersion 1.0.0
// @Title beego Test API
// @Description
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/b0ralgin/test-beego/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/signup", &controllers.LoginController{}, "post:SignUp"),
		beego.NSRouter("/signin", &controllers.LoginController{}, "post:SignIn"),
		//	beego.NSRouter("/reset", &controllers.LoginController{}, "post:Reset"),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.Router("/", &controllers.UserController{})
}
