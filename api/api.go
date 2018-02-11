package api

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

func familyProfileHandler(ctx iris.Context) {

}

func loginHandler(ctx iris.Context) {
	// TODO
	username := ctx.PostValue("username")
	_ = ctx.PostValue("password")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": username,
		"ur":  "u",
		"in":  time.Now().Unix(),
	})

	s, err := token.SignedString([]byte("foobar"))
	if err != nil {
		panic(err)
	}

	ctx.Header(jwtMiddleware.DefaultContextKey, fmt.Sprintf("bearer %v", s))
}

type familyApi struct {
	privateKey string
}

func NewApi() *familyApi {
	return &familyApi{
		privateKey: "foobar",
	}
}

func (api *familyApi) Run() {
	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	m := prometheusMiddleware.New("family", 300, 1200, 5000)

	app.Use(m.ServeHTTP)

	jwtHandler := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(api.privateKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	app.PartyFunc("/v1", func(v1 iris.Party) {
		v1.Use(jwtHandler.Serve)
		v1.PartyFunc("/users", func(users iris.Party) {
			users.Get("/{id:string}/profile", familyProfileHandler)
		})
	})
	app.Get("/metrics", iris.FromStd(prometheus.Handler()))

	app.Post("/login", loginHandler)

	app.Run(iris.Addr(":1527"))
}
