package todos

import (
	"context"
	"fmt"
	"go_api/database"
	"go_api/middleware/auth"
	"go_api/model"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginVals model.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		fmt.Println(jwt.ErrMissingLoginValues)
	}
	conn := database.Open()
	defer conn.Close(context.Background())
	var sql string

	sql = `SELECT * FROM public.user WHERE username=$1 AND password=$2`

	var person model.Person
	if err := conn.QueryRow(context.Background(), sql, loginVals.Username, loginVals.Password).Scan(&person.Id, &person.Username, &person.Password, &person.Age, &person.Role); err != nil {
		fmt.Println(err)
	} else {
		token, err := auth.GenerToken(person.Username, person.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"role":  person.Role,
		})
	}

}

func All(c *gin.Context) {

	conn := database.Open()
	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), `SELECT * FROM public.user`)

	for rows.Next() {
		var user model.Person
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Age, &user.Role); err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, user)
	}

}

func Add(c *gin.Context) {
	conn := database.Open()
	defer conn.Close(context.Background())

	var person model.Person
	if err := c.BindJSON(&person); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if _, err := conn.Exec(context.Background(), `insert into public.user(username, password, age) values($1, $2, $3)`, person.Username, person.Password, person.Age); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Panicln(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}

func FetchOne(c *gin.Context) {
	id := c.Param("id")

	conn := database.Open()
	defer conn.Close(context.Background())

	var person model.Person
	if err := conn.QueryRow(context.Background(), `SELECT * FROM public.user WHERE id=$1`, id).Scan(&person.Id, &person.Username, &person.Password, &person.Age, &person.Role); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, person)
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

	if _, err := conn.Exec(context.Background(), `update public.user set username=$1, age=$2, password=$3, role=$4 where id=$5`, putUser.Username, putUser.Age, putUser.Password, putUser.Role, id); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "更新成功"})

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

func AddCasbin(c *gin.Context) {
	var casbind model.CasbinBind
	if err := c.BindJSON(&casbind); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	casbin := auth.CasbinModel{
		Ptype:    casbind.Ptype,
		Rolename: casbind.Rolename,
		Path:     casbind.Path,
		Method:   casbind.Method,
	}
	var casbins = auth.CasbinModel{}

	isok, err := casbins.AddCasbin(casbin)
	if isok {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "保存成功",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     err,
		})
	}
}
