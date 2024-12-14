package main

import (
	"fmt"

	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	host := os.Getenv("DB_HOST")     // ? or the Docker service name if running in another container
	port := 5432                     // ? default PostgreSQL port
	user := os.Getenv("DB_USER")     // ? as defined in docker-compose.yml
	password := os.Getenv("DB_PASS") // ? as defined in docker-compose.yml
	dbname := os.Getenv("DB_NAME")   // ? as defined in docker-compose.yml

	// ? Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// ? New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // ? io writer
		logger.Config{
			SlowThreshold: time.Second, // ? Slow SQL threshold
			LogLevel:      logger.Info, // ? Log level
			Colorful:      true,        // ? Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // ? add Logger
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// ? Migrate the schema
	db.AutoMigrate(&Book{})
	fmt.Println("Database migration completed!")

	/*
		? Create
		newBook := &Book{
			Name:        "kate",
			Author:      "katetest",
			Description: "test_kate",
			Price:       4500,
		}

		CreateBook(db, newBook)
	*/

	/*
		? FindOne
		currentBook := GetBook(db, 5)

		currentBook.Name = "topUpdate"
		currentBook.Price = 500

		? Update
		UpdateBook(db, currentBook)
		fmt.Println(currentBook)
	*/

	/*
		? Delete
		DeleteBook(db, 5)
	*/

	/*
		? criteria books
		searchBook := SearchBook(db, "kate")

		for _, book := range searchBook {
			fmt.Println(book.ID, book.Name, book.Author, book.Price)
		}
	*/
}
