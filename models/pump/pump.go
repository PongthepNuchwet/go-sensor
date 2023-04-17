package pump

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PongthepNuchwet/go-sensor/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pump struct {
	ID        uint   `gorm:"primary key;autoIncrement" json:"id"`
	Value     string `json:"Value"`
	CreatedAt time.Time
}

func MigratePump(db *gorm.DB) error {
	err := db.AutoMigrate(&Pump{})
	return err
}

func CreatePump(context *fiber.Ctx) error {
	pump := Pump{}

	err := context.BodyParser(&pump)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = database.DBConn.Create(&pump).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create Pump"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pump has been added"})
	return nil
}

func GetPumps(context *fiber.Ctx) error {
	pumpModels := &[]Pump{}

	err := database.DBConn.Find(pumpModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get pump"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pump fetched successfully",
		"data":    pumpModels,
	})
	return nil
}

const YYYYMMDD = "2006-01-02"

func GetPumpsOntheday(context *fiber.Ctx) error {
	pumpModels := &[]Pump{}

	date := context.Params("date")
	fmt.Println("the date is", date)

	start, err := time.Parse(YYYYMMDD, date)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "invalid date"})
		return err
	}
	year, month, day := start.Date()
	startStr := fmt.Sprintf("%d-%s-%d", year, month, day)
	end := start.AddDate(0, 0, 1)
	year, month, day = end.Date()
	endStr := fmt.Sprintf("%d-%s-%d", year, month, day)
	fmt.Println("start", start)
	fmt.Println("startStr", startStr)
	fmt.Println("end", end)
	fmt.Println("endStr", endStr)

	err = database.DBConn.Where("created_at BETWEEN ? AND ?", startStr, endStr).Find(pumpModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get pump"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pump fetched successfully",
		"data":    pumpModels,
	})
	return nil
}
