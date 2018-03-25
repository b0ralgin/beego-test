package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/b0ralgin/test-beego/models"
	_ "github.com/b0ralgin/test-beego/routers"
	"github.com/b0ralgin/test-beego/utilities"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	err := utilities.MongoStartup()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestGet is a sample to run an endpoint test

func TestSignUp(t *testing.T) {
	user := models.User{
		Username: "test",
		Password: "test",
	}
	mUser, _ := json.Marshal(user)
	reqReader := bytes.NewReader(mUser)
	r, _ := http.NewRequest("POST", "/v1/signup", reqReader)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())
	response := models.User{}
	json.Unmarshal(w.Body.Bytes(), &response)
	Convey("Subject: SingUp Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("The result Should Be Equal", func() {
			So(response.Username, ShouldEqual, user.Username)
		})
	})
}

func TestSignIn(t *testing.T) {
	user := models.User{
		Username: "test",
		Password: "test",
	}
	mUser, _ := json.Marshal(user)
	reqReader := bytes.NewReader(mUser)
	r, _ := http.NewRequest("POST", "/v1/signin", reqReader)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())
	var token struct {
		Token string
	}
	json.Unmarshal(w.Body.Bytes(), &token)

	Convey("Subject: SingUp Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("The result Should Be Equal", func() {
			So(token.Token, ShouldNotBeNil)
		})
	})
}
