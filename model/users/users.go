package users

import "gorm.io/gorm"

type Base struct {
	gorm.Model
}
type Person struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Age      int    `json:"age"`
}

type Persons struct {
	Persons []Person `json:"persons"`
}
type AddUser struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Age      int    `json:"age"`
}

type UpdateUser struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Age      int    `json:"age"`
}

type User struct {
	UserName string
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type LoginDB struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Postgresql struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

type CasbinBind struct {
	Ptype    string `json:"ptype"`
	Rolename string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}
