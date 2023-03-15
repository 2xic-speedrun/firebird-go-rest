package main

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/nakagami/firebirdsql"
)

type SqlParams struct {
	Sql    string        `json:"sql"`
	Params []interface{} `json:"params"`
}

func rawExecuteSql(
	params SqlParams,
) map[string][]interface{} {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	USERNAME := os.Getenv("USERNAME")
	PASSWORD := os.Getenv("PASSWORD")
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	DATABASE := os.Getenv("DATABASE")
	CONNECTION_OPTIONS := os.Getenv("CONNECTION_OPTIONS")

	connectionString := strings.Join([]string{
		USERNAME, ":", PASSWORD, "@", HOST, ":", PORT, DATABASE,
		CONNECTION_OPTIONS,
	}, "")

	var db *sql.DB
	db, err = sql.Open("firebirdsql", connectionString)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(params.Sql, params.Params...)
	if err != nil {
		panic(err)
	}
	row_json, err := SQLToJSON(rows)
	if err != nil {
		panic(err)
	}
	return row_json
}

func executeSql(c *gin.Context) {
	var params SqlParams

	if err := c.BindJSON(&params); err != nil {
		return
	}

	row_json := rawExecuteSql(params)

	c.JSON(200, row_json)
}

func sayHello(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func main() {
	router := gin.New()
	router.GET("/", sayHello)
	router.POST("/", executeSql)

	router.Run("0.0.0.0:8080")
}
