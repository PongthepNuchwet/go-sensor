package main

import (
	"log"
	"os"

	"github.com/PongthepNuchwet/go-sensor/book"
	"github.com/PongthepNuchwet/go-sensor/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) HelloWorld(c *fiber.Ctx) error {
	text := []byte("Hello, World!")
	c.Send(text)
	return nil
}
func (r *Repository) SetupRoutes(app *fiber.App) {

	app.Get("/", r.HelloWorld)

	api := app.Group("/api")
	api.Post("/create_books", book.CreateBook)
	api.Delete("delete_book/:id", book.DeleteBook)
	api.Get("/get_books/:id", book.GetBookByID)
	api.Get("/books", book.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = book.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
