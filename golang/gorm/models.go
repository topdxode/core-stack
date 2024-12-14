package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Price       uint
}

func CreateBook(db *gorm.DB, book *Book) {
	result := db.Create(book)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}
	fmt.Println("Book created successfully")
}

func GetBook(db *gorm.DB, id uint) *Book {
	var book Book

	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return &book
}

func UpdateBook(db *gorm.DB, book *Book) {
	result := db.Save(book)
	if result.Error != nil {
		log.Fatalf("Error updating book: %v", result.Error)
	}
	fmt.Println("Book updated successfully")
}

func DeleteBook(db *gorm.DB, id uint) {
	var book Book

	result := db.Delete(&book, id)
	// ? use Unscoped() == delete permanant;
	// ! result := db.Unscoped().Delete(&book, id)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
	}
	fmt.Println("Book deleted successfully")
}

func SearchBook(db *gorm.DB, bookName string) []Book {
	var books []Book
	// ! slice [] standard send address

	result := db.Where("name = ?", bookName).Order("price DESC").Find(&books)
	if result.Error != nil {
		log.Fatalf("Search book failed: %v", result.Error)
	}
	return books
}
