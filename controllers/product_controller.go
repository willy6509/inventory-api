package controllers

import (
	"inventory-api/config"
	"inventory-api/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	if result := config.DB.Create(&product); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status": "success", 
		"data": product,
	})
}

func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product
	config.DB.Find(&products)
	return c.JSON(fiber.Map{"status": "success", "data": products})
}