package main

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title       string
	Description string
	ImagePath   string
	Author      Author
	AuthorId    uint
}

type Author struct {
	gorm.Model
	AuthorName string
	Password   string
	Blog       []Blog `json:"-"`
}
