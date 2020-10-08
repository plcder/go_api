package main

import (
	"go_api/controller/todos"
	"go_api/middleware/auth"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	authorized := r.Group("/")
	authorized.Use(auth.AuthRequired)
	authorized.Use(auth.AuthCheckRole())

	{
		r.POST("/add", todos.Add)
		authorized.GET("/api/users", todos.All)
		authorized.POST("/api/add", todos.Add)
		authorized.PUT("/api/update/:id", todos.Update)
		authorized.DELETE("/api/delete/:id", todos.Delete)
		authorized.POST("/api/todo/:id", todos.FetchOne)
		authorized.POST("api/addcabin", todos.AddCasbin)

	}
	r.POST("/api/login", todos.Login)

	r.Run()

}
