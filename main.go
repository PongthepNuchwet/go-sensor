package main

import (
	"log"
	"os"

	"github.com/PongthepNuchwet/go-sensor/database"
	"github.com/PongthepNuchwet/go-sensor/models/books"
	"github.com/PongthepNuchwet/go-sensor/models/pump"
	"github.com/PongthepNuchwet/go-sensor/models/waterFlow"
	"github.com/PongthepNuchwet/go-sensor/models/waterPressure"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func HelloWorld(c *fiber.Ctx) error {
	text := []byte("Hello, World!")
	c.Send(text)
	return nil
}
func SetupRoutes(app *fiber.App) {

	app.Get("/", HelloWorld)

	api := app.Group("/api")
	api.Post("/create_books", books.CreateBook)
	api.Delete("/delete_book/:id", books.DeleteBook)
	api.Get("/get_books/:id", books.GetBookByID)
	api.Get("/books", books.GetBooks)

	api.Post("/pumps", pump.CreatePump)
	api.Get("/pumps", pump.GetPumps)
	api.Get("/pumps/:date", pump.GetPumpsOntheday)

	api.Post("/waterflow", waterFlow.CreateWaterFlow)
	api.Get("/waterflow", waterFlow.GetWaterFlow)
	api.Get("/waterflow/:date", waterFlow.GetWaterFlowOntheday)

	api.Post("/waterPressure", waterPressure.CreateWaterPressure)
	api.Get("/waterPressure", waterPressure.GetWaterPressure)
	api.Get("/waterPressure/:date", waterPressure.GetWaterPressureOntheday)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}

	err = books.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	err = pump.MigratePump(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	err = waterFlow.MigratewaterFlow(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	err = waterPressure.MigrateWaterPressure(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	database.DBConn = db

	app := fiber.New()
	app.Use(cors.New())

	SetupRoutes(app)
	app.Listen(":8080")
}
