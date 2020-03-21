package todos

import (
	"context"
	"fmt"
	"go_api/database"
	"go_api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func All(c *gin.Context) {

	conn := database.Open()
	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), `SELECT * FROM users`)

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, user)
	}

}

func Add(c *gin.Context) {
	conn := database.Open()
	defer conn.Close(context.Background())

	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if _, err := conn.Exec(context.Background(), `insert into users(name, age) values($1, $2)`, user.Name, user.Age); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
func Update(c *gin.Context) {
	id := c.Param("id")

	var putUser model.UpdateUser

	if err := c.BindJSON(&putUser); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	conn := database.Open()
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), `update users set name=$1, age=$2 where id=$3`, putUser.Name, putUser.Age, id); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})

}

func Delete(c *gin.Context) {
	id := c.Param("id")

	conn := database.Open()
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), `delete from users where id=$1`, id); err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
