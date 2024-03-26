package main

import (
	"new_projects/controllers"
	"new_projects/initializers"
	"new_projects/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.POST("/create", controllers.CreateUser)
	r.POST("/login", controllers.Login)

	r.Use(middlewares.Authentication(), middlewares.Authoriztion())

	r.GET("/get", controllers.GetUser)
	r.GET("/get/:id", controllers.GetUserByID)
	r.PUT("/update/:id", controllers.UpdateUser)
	r.DELETE("/delete/:id", controllers.DeleteUser)
	r.Run()
}
