package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/next-bytes/syo-back/internal/database"
	"github.com/next-bytes/syo-back/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewPost struct {
	Author  string
	Content string
	Topics  []string
}

func GetPosts(c *fiber.Ctx) error {
	cursor, err := database.PostsCollection.Find(database.Ctx, bson.D{{Key: "answer", Value: nil}}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	posts := &[]models.Post{}
	if err = cursor.All(database.Ctx, posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(*posts)
}

func CreatePost(c *fiber.Ctx) error {

	body := &NewPost{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(body.Author) < 3 || len(body.Content) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author or content of message",
		})
	}

	var topics []string
	if len(body.Topics) == 0 {
		topics = []string{}
	} else {
		topics = body.Topics
	}

	newPost := &models.Post{
		ID:        primitive.NewObjectID(),
		Author:    body.Author,
		Content:   body.Content,
		Answer:    nil,
		Data:      models.PostData{},
		Topics:    topics,
		CreatedAt: time.Now(),
	}
	database.PostsCollection.InsertOne(database.Ctx, *newPost)

	return c.Status(fiber.StatusCreated).JSON(newPost)
}

func GetPostById(c *fiber.Ctx) error {
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
	post.Data.Views++
	database.PostsCollection.UpdateOne(database.Ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "data", Value: post.Data}}}})
	return c.Status(fiber.StatusOK).JSON(post)
}
