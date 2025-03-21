package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BlogWithOutAuthor struct {
	ID          uint
	Title       string
	Description string
	ImagePath   string
	AuthorId    uint
}

func getBlogs(db *gorm.DB, c *fiber.Ctx) error {
	var blogs []BlogWithOutAuthor
	result := db.Model(&Blog{}).Find(&blogs)
	if result.Error != nil {
		log.Fatalf("Error finding all blog: %v", result.Error)
	}
	return c.JSON(blogs)
}

func getBlogsWithAuthor(db *gorm.DB, c *fiber.Ctx) error {
	var blogs []Blog
	result := db.Preload("Author").Find(&blogs)
	if result.Error != nil {
		log.Fatalf("Error finding all blog: %v", result.Error)
	}
	return c.JSON(blogs)
}
