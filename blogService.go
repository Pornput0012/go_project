package main

import (
	"fmt"
	"log"
	"strconv"

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

func getBlog(db *gorm.DB, c *fiber.Ctx) error {
	blogIdStr := c.Params("blog_id")
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var blog BlogWithOutAuthor
	if result := db.Model(&Blog{}).First(&blog, blogId); result.Error != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(blog)
}

func getBlogWithAuthor(db *gorm.DB, c *fiber.Ctx) error {
	blogIdStr := c.Params("blog_id")
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var blog Blog
	if result := db.Preload("Author").First(&blog, blogId); result.Error != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(blog)
}

func createBlog(db *gorm.DB, c *fiber.Ctx) error {
	// รับไฟล์ภาพจากฟอร์ม
	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to retrieve image file",
		})
	}

	// สร้างพาธสำหรับบันทึกไฟล์ภาพ
	imagePath := fmt.Sprintf("images/%s", image.Filename)

	// บันทึกไฟล์ลงเซิร์ฟเวอร์
	if err := c.SaveFile(image, imagePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image file",
		})
	}

	// แปลง Body เป็น struct
	blog := new(BlogWithOutAuthor)
	if err := c.BodyParser(blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// กำหนด ImagePath ที่ได้จากการอัปโหลด
	blog.ImagePath = imagePath

	// สร้าง Blog ในฐานข้อมูล
	newBlog := Blog{
		Title:       blog.Title,
		Description: blog.Description,
		ImagePath:   blog.ImagePath,
		AuthorId:    blog.AuthorId,
	}

	if result := db.Create(&newBlog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to create blog",
			"errorMsg": result.Error,
		})
	}

	// ส่งข้อมูลกลับ
	return c.Status(fiber.StatusCreated).JSON(newBlog)
}
