package main

import (
    "github.com/gin-gonic/gin"
    "os"
    "database/sql"
    _ "github.com/nakagami/firebirdsql"
    "github.com/joho/godotenv"
    "strings"
)

type SqlParams struct {
    Sql     string   `json:"sql"`
    Params  []interface{} `json:"params"`
}

func executeSql(c *gin.Context) {
    var params SqlParams

    if err := c.BindJSON(&params); err != nil {
        return
    }

    err := godotenv.Load()
    if err != nil {
      panic("Error loading .env file")
    }

    USERNAME := os.Getenv("USERNAME")
    PASSWORD := os.Getenv("PASSWORD")
    HOST := os.Getenv("HOST")
    PORT := os.Getenv("PORT")
    DATABASE := os.Getenv("DATABASE")

    connectionString := strings.Join([]string{
        USERNAME, ":", PASSWORD, "@", HOST, ":", PORT, "/", DATABASE,
    }, "")
  
    var db *sql.DB
    db, _ = sql.Open("firebirdsql", connectionString)
    rows, err := db.Query(params.Sql, params.Params...)
    if err != nil {
        panic(err);
    }
    row_json, _ := SQLToJSON(rows)
    c.JSON(200, row_json)
}

func main() {
//    router := gin.Default()
    router := gin.New()
    router.POST("/", executeSql)

//    router.Run("localhost:8080")
    router.Run("0.0.0.0:8080")
}
