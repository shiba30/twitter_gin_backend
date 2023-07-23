/*
	全てのルーティングを行う
*/

package api

import (
	"example.com/golang_twitter/api/user"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	api := router.Group("/api")
	{
		user.Routes(api)
	}
}
