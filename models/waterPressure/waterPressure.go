package waterPressure

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PongthepNuchwet/go-sensor/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WaterPressure struct {
	ID        uint `gorm:"primary key;autoIncrement" json:"id"`
	Value     int  `json:"Value"`
	CreatedAt time.Time
}

func MigrateWaterPressure(db *gorm.DB) error {
	err := db.AutoMigrate(&WaterPressure{})
	return err
}

func CreateWaterPressure(context *fiber.Ctx) error {
	waterPressure := WaterPressure{}

	err := context.BodyParser(&waterPressure)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = database.DBConn.Create(&waterPressure).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create WaterPressure"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "WaterPressure has been added"})
	return nil
}

func GetWaterPressure(context *fiber.Ctx) error {
	waterPressureModel := &[]WaterPressure{}

	err := database.DBConn.Find(waterPressureModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get WaterPressure"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "WaterPressure fetched successfully",
		"data":    waterPressureModel,
	})
	return nil
}

const YYYYMMDD = "2006-01-02"

func GetWaterPressureOntheday(context *fiber.Ctx) error {
	waterPressureModel := &[]WaterPressure{}

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

	err = database.DBConn.Where("created_at BETWEEN ? AND ?", startStr, endStr).Find(waterPressureModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get WaterPressure"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "WaterPressure fetched successfully",
		"data":    waterPressureModel,
	})
	return nil
}
