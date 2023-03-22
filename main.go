package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	//"html/template"
	//"strings"
	globals "Todo_system/globals"
	middleware "Todo_system/middleware"
	routes "Todo_system/routes"
	database "Todo_system/database"
)


func main() {
	db := database.ConnectDB()



	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	routes.PublicRoutes(public, db)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private, db)

	router.Run("localhost:3000")
}