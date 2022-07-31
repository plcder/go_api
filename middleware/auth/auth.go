package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_api/database"

	"github.com/casbin/casbin/v2"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Account string `json:"username"`
	Role    string `json:"password"`
	jwt.StandardClaims
}
type CasbinModel struct {
	ID       int    `json:"id"`
	Ptype    string `json:"ptype"`
	Rolename string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

var jwtSecret = []byte("secret")

func GenerToken(username, role string) (string, error) {

	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	jwtId := username + strconv.FormatInt(nowTime.Unix(), 10)

	claims := Claims{
		Account: username,
		Role:    role,
		StandardClaims: jwt.StandardClaims{
			Audience:  username,
			ExpiresAt: expireTime.Unix(),
			Id:        jwtId,
			IssuedAt:  nowTime.Unix(),
			Issuer:    "go_api",
			NotBefore: nowTime.Add(10 * time.Second).Unix(),
			Subject:   username,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err

}

func AuthRequired(c *gin.Context) {
	token := c.GetHeader("Authorization")

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return jwtSecret, nil
		})
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		c.Abort()
		return
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {

		fmt.Println("account:", claims.Account)
		fmt.Println("role:", claims.Role)
		c.Set("account", claims.Account)
		c.Set("role", claims.Role)
		c.Next()
	} else {
		c.Abort()
		return
	}
}

func (c *CasbinModel) AddCasbin(cm CasbinModel) (bool, error) {
	e := Casbin()
	return e.AddPolicy(cm.Rolename, cm.Path, cm.Method)
}

func Casbin() *casbin.Enforcer {
	conf := database.Config()
	connectionString := "postgres://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	tableName := "casbin"
	adpter, err := casbinpgadapter.NewAdapter(db, tableName)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("config/auth_model.conf", adpter)
	if err != nil {
		panic(err)
	}
	return e
}

func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
			func(token *jwt.Token) (i interface{}, err error) {
				return jwtSecret, nil
			})
		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "token is malformed"
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "token could not be verified because of signing problems"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "signature validation failed"
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "token is expired"
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "token is not yet valid before sometime"
				} else {
					message = "can not handle this token"
				}
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		claims, _ := tokenClaims.Claims.(*Claims)

		role := claims.Role

		e := Casbin()

		res, err := e.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    "錯誤消息" + err.Error(),
			})
			c.Abort()
			return
		}
		if res {
			c.Next()
		} else {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"status": 301,
				"msg":    "很抱歉你沒有此權限",
			})
			c.Abort()
			return
		}
	}
}
