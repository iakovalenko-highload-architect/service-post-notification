package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"service-post-notification/cmd"
	customMiddlewares "service-post-notification/internal/middlewares"
	"service-post-notification/internal/usecase/token_manager"
	"service-post-notification/internal/utils/jwt"
	"service-post-notification/internal/websocket/handlers/post_feed_posted"
)

func main() {
	e := echo.New()

	rabbit := cmd.MustInitRabbit()
	defer rabbit.Close()

	tokenManager := token_manager.New(jwt.New(), cmd.MustInitTokenManagerConfig())

	postFeedPostedHandler := post_feed_posted.New(rabbit)

	middlewares := customMiddlewares.New(tokenManager)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.Auth)

	e.GET("/", postFeedPostedHandler.Handle)

	e.Logger.Fatal(e.Start(":1323"))
}
