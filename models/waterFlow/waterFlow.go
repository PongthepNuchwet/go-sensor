package waterFlow

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PongthepNuchwet/go-sensor/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WaterFlow struct {
	ID        uint `gorm:"primary key;autoIncrement" json:"id"`
	Value     int  `json:"Value"`
	CreatedAt time.Time
}

func MigratewaterFlow(db *gorm.DB) error {
	err := db.AutoMigrate(&WaterFlow{})
	return err
}

func CreateWaterFlow(context *fiber.Ctx) error {
	waterFlow := WaterFlow{}

	err := context.BodyParser(&waterFlow)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = database.DBConn.Create(&waterFlow).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create waterFlow"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "waterFlow has been added"})
	return nil
}

func GetWaterFlow(context *fiber.Ctx) error {
	waterFlowModel := &[]WaterFlow{}

	err := database.DBConn.Find(waterFlowModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get waterFlow"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "waterFlow fetched successfully",
		"data":    waterFlowModel,
	})
	return nil
}

const YYYYMMDD = "2006-01-02"

func GetWaterFlowOntheday(context *fiber.Ctx) error {
	waterFlowModel := &[]WaterFlow{}

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

	err = database.DBConn.Where("created_at BETWEEN ? AND ?", startStr, endStr).Find(waterFlowModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get waterFlow"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"data": waterFlowModel})
	return nil
}
