package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/next-bytes/syo-back/internal/database"
	"github.com/next-bytes/syo-back/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPostTopics(c *fiber.Ctx) error {
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Post ID",
		})
	}

	filter := bson.D{{Key: "_id", Value: postID}}
	post := &models.Post{}
	if err := database.PostsCollection.FindOne(database.Ctx, filter).Decode(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "This post not exists"})
	}

	return c.Status(fiber.StatusOK).JSON(post.Topics)
}

func GetTopicsRecommended(c *fiber.Ctx) error {
	cursor, err := database.PostsCollection.Find(database.Ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))

	posts := new([]models.Post)
	if err = cursor.All(database.Ctx, posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var topics []string = []string{}
	for _, post := range *posts {
		if len(topics) >= 5 {
			break
		}
		topics = append(topics, post.Topics...)
	}
	if len(topics) >= 5 {
		topics = append(topics[0:], topics[5:]...)
	}
	return c.Status(fiber.StatusOK).JSON(topics)
}
