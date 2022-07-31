package database

import (
	"go_api/struct/users"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgresql struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

func Open() *gorm.DB {

	conf := Config()
	var db *gorm.DB
	var err error

	dsn := "user=" + conf.User + " password=" + conf.Password + " dbname=" + conf.Name + " port=" + conf.Port + " sslmode=disable TimeZone=Asia/Taipei"
	if gin.Mode() == gin.ReleaseMode {

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	}

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&users.Person{})

	return db

}

func Config() (conf *Postgresql) {
	conf = new(Postgresql)
	yamlFile, _ := ioutil.ReadFile("database/Postgres.yaml")

	errUn := yaml.Unmarshal(yamlFile, conf)
	if errUn != nil {
		log.Fatalf("Unmarshal:", errUn)
	}
	return conf
}
