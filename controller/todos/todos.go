package todos

import (
	"fmt"
	"go_api/database"
	"go_api/middleware/auth"
	"go_api/model/users"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginVals users.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		fmt.Println(jwt.ErrMissingLoginValues)
	}
	db := database.Open()

	var sql string

	sql = `SELECT * FROM public.people WHERE username=? AND password=?`

	var person users.Person
	if err := db.Raw(sql, loginVals.Username, loginVals.Password).First(&person).Error; err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(person.Role)
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

	db := database.Open()
	var persons users.Persons

	rows := db.Raw(`SELECT * FROM public.user`).Scan(&persons)

	c.JSON(http.StatusOK, rows)
}

func Add(c *gin.Context) {
	db := database.Open()

	var person users.Person
	if err := c.BindJSON(&person); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err := db.Create(&person).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, person)

}

func FetchOne(c *gin.Context) {
	id := c.Param("id")

	db := database.Open()

	var person users.Person
	if err := db.Model(&person).Where("id = ?", id).First(&person).Error; err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, person)
}

func Update(c *gin.Context) {
	id := c.Param("id")

	var putUser users.UpdateUser

	if err := c.BindJSON(&putUser); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db := database.Open()

	if err := db.Model(&putUser).Where("id = ?", id).Error; err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "msg": "更新成功"})

}

func Delete(c *gin.Context) {
	id := c.Param("id")

	db := database.Open()
	var person users.Person

	if err := db.Model(&person).Where("id = ?", id).Error; err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func AddCasbin(c *gin.Context) {
	var casbind users.CasbinBind
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
