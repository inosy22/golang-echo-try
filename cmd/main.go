package main

import (
	"github.com/inosy22/golang-echo-try/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/redis", controller.RedisHello)
	e.POST("/redis", controller.RedisPost)
	e.PUT("/redis/:key", controller.RedisPut)
	e.GET("/redis/:key", controller.RedisGet)
	e.DELETE("/redis/:key", controller.RedisDelete)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello, World!")
}
