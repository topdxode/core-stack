package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func logMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	return c.Next()
}

func isAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	app := fiber.New()

	// ? Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // ? Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// ? Setup routes
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	app.Post("/login", login)

	app.Use(logMiddleware)

	// ? JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	bookGroup := app.Group("/book")

	bookGroup.Use(isAdmin)

	bookGroup.Get("/", getBooks)
	bookGroup.Get("/:id", getBook)
	bookGroup.Post("/", createBook)
	bookGroup.Put("/:id", updateBook)
	bookGroup.Delete("/:id", deleteBook)

	app.Post("/upload", uploadFile)

	// ? Use the environment variable for the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // ? Default port if not specified
	}

	app.Listen(":" + port)
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded successfully: " + file.Filename)
}

// ? Dummy user for example
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// ? Check credentials - In real world, you should check against a database
	if memberUser.Email != user.Email || memberUser.Password != user.Password {
		return fiber.ErrUnauthorized
	}

	// ? Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// ? Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// ? Generate encoded token
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   t,
	})
}

/*
	* use case struct dubplicate struct
	type Author stuct {
		Name 	string `json:"name"`
		Country string `json:"country"`
	}

	type Book stuct {
		Title string `json:"title"`
		Pages int 	 `json:"pages"`
	}

	type AuthorBookRequest struct {
		AuthorDetails 	Author  `json:"author"`
		BookDetails 	Book 	`json:"book`
	}

	? body example =>
	{
		"author": {
			"name": "George",
			"country": "United Kingdom"
		},
		"book": {
			"title": "hi bro!",
			"pages": 328
		}
	}
*/

/*
	* Template with Fiber (email, pdf)
	package main

	import (
		"github.com/gofiber/fiber/v2"
		"github.com/gofiber/template/html/v2"
	)

	func main() {
		? Initialize standard Go html template engine
		engine := html.New("./views", ".html")

		? Pass the engine to Fiber
		app := fiber.New(fiber.Config{
			Views: engine,
		})

		? Setup route
		app.Get("/", renderTemplate)

		app.Listen(":8080")
	}

	func renderTemplate(c *fiber.Ctx) error {
		? Render the template with variable data
		return c.Render("template", fiber.Map{
			"Name": "World",
		})
	}

	? views/template.html
	<!DOCTYPE html>
	<html>
	<head>
		<title>Example Template</title>
	</head>
	<body>
		<h1>Hello, {{.Name}}!</h1>
	</body>
	</html>
*/
