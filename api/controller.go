/*
	全てのルーティングを行う
*/

package api

import (
	"example.com/golang_twitter/api/user"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config config.Config
}

func NewController(cfg config.Config) *Controller {
	return &Controller{Config: cfg}
}

func (c *Controller) Routes(router *gin.Engine) {
	api := router.Group("/api")
	{
		user.SignupRoutes(api, c.Config)
		user.LoginRoutes(api, c.Config)
	}
}
