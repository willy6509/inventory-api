package controllers

import (
	"inventory-api/config"
	"inventory-api/models"
	"github.com/gofiber/fiber/v2"
)

// Lihat semua user yang statusnya PENDING
func GetPendingUsers(c *fiber.Ctx) error {
	var users []models.User
	config.DB.Where("status = ?", "pending").Find(&users)
	return c.JSON(users)
}

// Approve User (Ubah status jadi Active)
func ApproveUser(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Update status jadi active
	user.Status = "active"
	config.DB.Save(&user)

	return c.JSON(fiber.Map{
		"message": "User berhasil disetujui dan sekarang bisa login",
		"user":    user,
	})
}