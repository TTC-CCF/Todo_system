package controllers

import (
	// "fmt"

	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	database "Todo_system/database"
	globals "Todo_system/globals"
)

func RegisterGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"type":"is-danger is-light",
					"content": "Please logout first",
					"user":    user,
				})
			return
		}
		c.HTML(http.StatusOK, "register.html", gin.H{
			"type":"is-light", 
			"content": "Please enter username and password",
			"user":    user,
		})
	}
}

func RegisterPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c * gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first", "type":"is-danger is-light"})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")

		if database.EmptyUserPass(username, password) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Parameters can't be empty", "type":"is-danger is-light"})
			return
		}

		if database.CheckUserExist(db, username) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"content": "Username exist", "type":"is-danger is-light"} )
			return
		}

		if err := database.CreateUser(db, username, password); err != nil{
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{"content": "Internal server error, please try again", "type":"is-danger is-light"})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"type" : "is-ligth",
			"content": "Login with your new account again",
			"user":    user,
		})

	}
}

func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"type" : "is-danger is-light",
					"content": "Please logout first",
					"user":    user,
				})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"type" : "",
			"content": "",
			"user":    user,
		})
	}
}

func LoginPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first", "type":"is-danger is-light"})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")

		if database.EmptyUserPass(username, password) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Parameters can't be empty", "type":"is-danger is-light"})
			return
		}

		if !database.CheckUserPass(db, username, password) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password", "type":"is-danger is-light"})
			return
		}

		session.Set(globals.Userkey, username)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session", "type":"is-danger is-light"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		session.Delete(globals.Userkey)
		
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}
		
		c.Redirect(http.StatusTemporaryRedirect,  "/")
	}
}

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "hello world",
			"user":    user,
		})
	}
}

func DashboardGetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		todolist := database.GetTodoList(db, fmt.Sprint(user))
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"todolist": todolist,
			"content": "Create your Todoo here",
			"user":    user,
		})
	}
}

func DashboardPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)		
		method := c.Request.FormValue("method")
		if method == "createTodo"{
			title := c.PostForm("title")
			if err := database.CreateTodo(db, fmt.Sprint(user), title); err != nil{
				log.Println(err)
				return
			}
		}
		if method == "Done"{
			idstring := c.Request.PostForm.Get("ID")
			donestring := c.Request.PostForm.Get("Done")

			
			donebool, _ := strconv.ParseBool(donestring)
			idint, _:= strconv.Atoi(idstring)
			if err := database.DoneTodo(db, idint, donebool); err != nil{
				log.Println(err)
				return
			}
		}
		if method == "Delete"{
			idstring := c.Request.PostForm.Get("ID")
			idint, _:= strconv.Atoi(idstring)
			if err := database.DeleteTodo(db, idint); err != nil{
				log.Println(err)
				return
			}
		}
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}
