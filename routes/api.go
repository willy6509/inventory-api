package routes

import (
	"inventory-api/controllers"
	"inventory-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- AUTH (Google) ---
	api.Get("/auth/google", controllers.GoogleLogin)
	api.Get("/auth/google/callback", controllers.GoogleCallback)

	// --- ADMIN (Hanya bisa diakses User Active & Role Admin) ---
	admin := api.Group("/admin", middleware.Protected, middleware.AdminOnly)
	admin.Get("/pending-users", controllers.GetPendingUsers)
	admin.Put("/approve/:id", controllers.ApproveUser)

	// --- PRODUCTS (Hanya User Active) ---
	products := api.Group("/products", middleware.Protected)
	products.Get("/", controllers.GetAllProducts)
	products.Post("/", controllers.CreateProduct)

	// --- TRANSACTIONS (Hanya User Active) ---
	trx := api.Group("/transactions", middleware.Protected)
	trx.Post("/in", controllers.CreateTransactionIn)   // Barang Masuk
	trx.Get("/in", controllers.GetTransactionsIn)      // History Masuk
	trx.Post("/out", controllers.CreateTransactionOut) // Barang Keluar
	trx.Get("/out", controllers.GetTransactionsOut)    // History Keluar
}