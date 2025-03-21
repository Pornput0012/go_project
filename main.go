package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	app := fiber.New()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)
	dsn := "root:rootpassword@tcp(127.0.0.1:3310)/blogproject?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Can't Connect Database!")
	}

	db.AutoMigrate(&Blog{}, &Author{})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	blogGroupApi := app.Group("/blogs")
	blogGroupApi.Get("", func(c *fiber.Ctx) error {
		return getBlogs(db, c)
	})
	blogGroupApi.Get("/authors", func(c *fiber.Ctx) error {
		return getBlogsWithAuthor(db, c)
	})

	app.Listen("localhost:8080")

}
