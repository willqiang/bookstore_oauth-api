package app

import (
	"github.com/gin-gonic/gin"
	"github.com/willqiang/bookstore_oauth-api/src/clients/cassandra"
	"github.com/willqiang/bookstore_oauth-api/src/domain/access_token"
	"github.com/willqiang/bookstore_oauth-api/src/http"
	"github.com/willqiang/bookstore_oauth-api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session, err := cassandra.GetSession()
	if err != nil {
		panic(err)
	}
	session.Close()
	atService := access_token.NewService(db.NewRepository())
	atHandler := http.NewHandler(atService)
	//atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
