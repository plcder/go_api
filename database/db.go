package database

import (
	"context"
	"io/ioutil"
	"log"

	"go_api/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	yaml "gopkg.in/yaml.v2"
)

func Open() *pgx.Conn {

	conf := Config()
	var conn *pgx.Conn
	var err error

	if gin.Mode() == gin.ReleaseMode {

		conn, err = pgx.Connect(context.Background(), "postgres://"+conf.User+":"+conf.Password+"@"+conf.Host+":"+conf.Port+"/postgres")
	} else {
		conn, err = pgx.Connect(context.Background(), "postgres://"+conf.User+":"+conf.Password+"@"+conf.Host+":"+conf.Port+"/postgres")
	}

	if err != nil {
		panic(err)
	}

	return conn

}

func Config() (conf *model.Postgresql) {
	conf = new(model.Postgresql)
	yamlFile, _ := ioutil.ReadFile("database/Postgres.yaml")

	errUn := yaml.Unmarshal(yamlFile, conf)
	if errUn != nil {
		log.Fatalf("Unmarshal:", errUn)
	}
	return conf
}
