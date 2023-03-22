package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	controllers "Todo_system/controllers"
)

func PublicRoutes(g *gin.RouterGroup, db *gorm.DB) {


	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler(db))
	g.GET("/register", controllers.RegisterGetHandler())
	g.POST("/register", controllers.RegisterPostHandler(db))
	g.GET("/", controllers.IndexGetHandler())

}

func PrivateRoutes(g *gin.RouterGroup, db *gorm.DB) {

	g.GET("/dashboard", controllers.DashboardGetHandler())
	g.GET("/logout", controllers.LogoutGetHandler())
}