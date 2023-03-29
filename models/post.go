package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Author    string             `json:"author" bson:"author"`
	Content   string             `json:"content" bson:"content"`
	Answer    interface{}        `json:"answer" bson:"answer"`
	Data      PostData           `json:"data" bson:"data"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type PostData struct {
	Likes    int `json:"likes" bson:"likes"`
	Views    int `json:"views" bson:"views"`
	Comments int `json:"comments" bson:"comments"`
}
