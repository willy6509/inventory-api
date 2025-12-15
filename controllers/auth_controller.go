package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"inventory-api/config"
	"inventory-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Redirect ke Google
func GoogleLogin(c *fiber.Ctx) error {
	url := config.GetGoogleConfig().AuthCodeURL("state-token")
	return c.Redirect(url)
}

// Callback dari Google
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	googleConfig := config.GetGoogleConfig()

	// 1. Tukar Code dengan Token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal koneksi ke Google"})
	}

	// 2. Ambil Info User
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data user"})
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode respon Google"})
	}

	// 3. Cek Database
	var user models.User
	result := config.DB.Where("email = ?", googleUser.Email).First(&user)

	// Jika User Belum Ada -> Register PENDING
	if result.RowsAffected == 0 {
		newUser := models.User{
			GoogleID: googleUser.ID,
			Email:    googleUser.Email,
			FullName: googleUser.Name,
			Avatar:   googleUser.Picture,
			Role:     "staff",
			Status:   "pending", // Default Pending
		}
		config.DB.Create(&newUser)
		return c.Status(403).JSON(fiber.Map{
			"message": "Pendaftaran berhasil. Akun Anda berstatus PENDING. Hubungi Admin untuk aktivasi.",
		})
	}

	// Jika User Ada tapi PENDING
	if user.Status == "pending" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akun Anda belum diaktifkan oleh Admin.",
		})
	}

	// Jika User Ada & ACTIVE -> Buat Token JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login Berhasil",
		"token":   tokenString,
		"user":    user,
	})
}