package main

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/nakagami/firebirdsql"
	"github.com/pquerna/ffjson/ffjson"
)

type SqlQuery struct {
	Sql    string        `json:"sql"`
	Params []interface{} `json:"params"`
}

func connectToDb() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	USERNAME := os.Getenv("FIREBIRD_REST_USERNAME")
	PASSWORD := os.Getenv("FIREBIRD_REST_PASSWORD")
	HOST := os.Getenv("FIREBIRD_REST_HOST")
	PORT := os.Getenv("FIREBIRD_REST_PORT")
	DATABASE := os.Getenv("FIREBIRD_REST_DATABASE")
	CONNECTION_OPTIONS := os.Getenv("FIREBIRD_REST_CONNECTION_OPTIONS")

	connectionString := strings.Join([]string{
		USERNAME, ":", PASSWORD, "@", HOST, ":", PORT, DATABASE,
		CONNECTION_OPTIONS,
	}, "")

	var db *sql.DB
	db, err = sql.Open("firebirdsql", connectionString)
	if err != nil {
		panic(err)
	}
	return db
}

func rawExecuteSql(
	db *sql.DB,
	query SqlQuery,
) map[string][]interface{} {
	rows, err := db.Query(query.Sql, query.Params...)
	if err != nil {
		panic(err)
	}

	row_json, err := SQLToJSON(rows)
	if err != nil {
		panic(err)
	}
	return row_json
}

func executeSql(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var params SqlQuery

		if err := c.BindJSON(&params); err != nil {
			panic(err)
		}

		row_json := rawExecuteSql(db, params)

		data, err := ffjson.Marshal(row_json)
		if err != nil {
			panic(err)
		}

		c.String(200, string(data))
	}

	return gin.HandlerFunc(fn)
}

func sayHello(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	db := connectToDb()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/", sayHello)
	router.POST("/", executeSql(db))

	router.Run("0.0.0.0:8080")
}
