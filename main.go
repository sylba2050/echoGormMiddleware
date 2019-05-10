//  https://medium.com/veltra-engineering/echo-middleware-in-golang-90e1d301eb27
package main

import (
    "net/http"

    "github.com/labstack/echo"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := gorm.Open("sqlite3", "main.sqlite3")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    db.AutoMigrate(&auth{})

    e := echo.New()
    g := e.Group("", customMiddleware(db))
    g.GET("/", test(db))

    e.Start(":8080")
}

func test(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.NoContent(http.StatusOK)
    }
}

type auth struct {
    gorm.Model
    UserId string `json:"userid" form:"userid" query:"userid"`
    PW string `json:"pw" form:"pw" query:"pw"`
}

func customMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            a := new(auth)
            a.UserId = "mstn_"
            a.PW = "admin"
            db.Create(&a)

            err := next(c)
            return err
        }
    }
}
