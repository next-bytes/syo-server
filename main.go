package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/next-bytes/syo-back/internal/database"
	"github.com/next-bytes/syo-back/routes"
)

func registerRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	{
		v1Posts := v1.Group("/posts")
		v1Posts.Get("/", routes.GetPosts)
		v1Posts.Post("/", routes.CreatePost)
		v1Posts.Get("/:id", routes.GetPostById)

		v1Posts.Get("/:id/comments", routes.GetPostComments)
		v1Posts.Post("/:id/comments", routes.CreatePostComment)

		v1Posts.Get("/:id/topics", routes.GetPostTopics)

		v1.Get("/topics", routes.GetTopicsRecommended)
	}
}

func startServer() {

	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}

	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println(c.Method(), c.Path(), "-", c.IP(), "|", c.Response().StatusCode())
		return c.Next()
	})

	registerRoutes(app)

	ip := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	app.Listen(ip)
}

func main() {
	startServer()
}
