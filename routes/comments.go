package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/next-bytes/syo-back/internal/database"
	"github.com/next-bytes/syo-back/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPostComments(c *fiber.Ctx) error {
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Post ID",
		})
	}

	filter := bson.D{{Key: "answer", Value: postID}}
	cursor, err := database.PostsCollection.Find(database.Ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	post := &[]models.Post{}
	if err = cursor.All(database.Ctx, post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(post)
}

type NewComment struct {
	Author  string
	Content string
	Answer  string
}

func CreatePostComment(c *fiber.Ctx) error {
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Post ID",
		})
	}
	body := new(NewComment)
	if err = c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(body.Author) < 3 || len(body.Content) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author or content of message",
		})
	}

	filter := bson.D{{Key: "_id", Value: postID}}
	originalPost := new(models.Post)
	if err = database.PostsCollection.FindOne(database.Ctx, filter).Decode(originalPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "This post not exists"})
	}

	originalPost.Data.Comments++
	database.PostsCollection.UpdateOne(database.Ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "data", Value: originalPost.Data}}}})

	newPost := &models.Post{
		ID:        primitive.NewObjectID(),
		Author:    body.Author,
		Content:   body.Content,
		Answer:    originalPost.ID,
		Data:      models.PostData{},
		Topics:    originalPost.Topics,
		CreatedAt: time.Now(),
	}

	database.PostsCollection.InsertOne(database.Ctx, *newPost)

	return c.Status(fiber.StatusCreated).JSON(newPost)
}
