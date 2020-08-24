package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type GinRouter struct {
	Router *gin.Engine
}

var ginEngine = &GinRouter{Router: nil}

func CreateRouter() *GinRouter {
	if ginEngine.Router == nil {
		fmt.Println("Creating new router ...")
		ginEngine.Router = gin.Default()
		ginEngine.Router.RedirectTrailingSlash = false
		ginEngine.Router.RemoveExtraSlash = true
	}
	return ginEngine
}

func (r GinRouter) Run(certFile, certKey string) {
	err := r.Router.RunTLS(":8443", certFile, certKey)
	if err != nil {
		log.Fatal(err)
	}
}

func (r GinRouter) AddHandler(handler func(GinRouter)) {
	handler(r)
}
