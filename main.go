package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type UserCredential struct {
	Id           uint   `db:"id"`
	UserName     string `db:"user_name"`
	IsBlocked    string `db:"is_blocked"`
	UserPassword string `db:"user_password"`
}
type Login struct {
	User     string `json:"userName"  binding:"required"`
	Password string `json:"userPassword" binding:"required"`
}

func main() {
	apiHost := "localhost"
	apiPort := "8888"
	dbHost := "159.223.42.164"
	dbPort := "5432"
	dbName := "enigma"
	dbUser := "postgres"
	dbPassword := "P@ssw0rd"
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	conn, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer func(conn *sqlx.DB) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
	r := gin.Default()
	route := r.Group("/enigma")
	route.POST("/auth", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var userCred = UserCredential{}
		sql := "SELECT * FROM user_credentials where user_name=$1 and user_password=$2"
		fmt.Println(sql)

		err := conn.Get(&userCred, sql, login.User, login.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	listenAddress := fmt.Sprintf("%s:%s", apiHost, apiPort)
	err = r.Run(listenAddress)
	if err != nil {
		panic(err)
	}
}

/*
Sample Body Request
{
    "userName":"myuser' or 'foo' = 'foo",
    "userPassword":"myuser' or 'foo' = 'foo"
}
*/
