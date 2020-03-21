package main

import (
	"go_api/controller/todos"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/api/todos", todos.All)
	r.POST("/api/add", todos.Add)
	r.PUT("/api/update/:id", todos.Update)
	r.DELETE("/api/delete/:id", todos.Delete)
	r.Run()

}
