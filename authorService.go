package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func register(db *gorm.DB, c *fiber.Ctx) error {
	author := new(Author)
	if err := c.BodyParser(author); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if result := db.Create(&author); result.Error != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(author)
}

func getAuthors(db *gorm.DB, c *fiber.Ctx) error {
	var authors []Author

	if result := db.Find(&authors); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to find authors",
		})
	}

	return c.JSON(authors)
}
