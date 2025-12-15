package controllers

import (
	"inventory-api/config"
	"inventory-api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// --- BARANG MASUK (IN) ---

// 1. Catat Barang Masuk (Otomatis Tambah Stok)
func CreateTransactionIn(c *fiber.Ctx) error {
	txData := new(models.TransactionIn)

	// Parsing input JSON
	if err := c.BodyParser(txData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input data tidak valid"})
	}

	// Gunakan Database Transaction (Agar aman: jika satu gagal, semua batal)
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// A. Simpan data transaksi ke tabel transactions_in
		if err := tx.Create(&txData).Error; err != nil {
			return err
		}

		// B. Update stok di tabel products (Stok + Quantity)
		// Kita menggunakan gorm.Expr agar update bersifat atomic di level database
		if err := tx.Model(&models.Product{}).
			Where("id = ?", txData.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", txData.Quantity)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses transaksi masuk: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Barang masuk berhasil dicatat & stok bertambah",
		"data":    txData,
	})
}

// 2. Lihat Riwayat Barang Masuk
func GetTransactionsIn(c *fiber.Ctx) error {
	var transactions []models.TransactionIn
	// Preload("Product") akan otomatis mengambil detail nama barangnya juga
	config.DB.Preload("Product").Order("date desc").Find(&transactions)
	return c.JSON(transactions)
}

// --- BARANG KELUAR (OUT) ---

// 3. Catat Barang Keluar (Otomatis Kurang Stok)
func CreateTransactionOut(c *fiber.Ctx) error {
	txData := new(models.TransactionOut)

	if err := c.BodyParser(txData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input data tidak valid"})
	}

	// Cek Stok Dulu (Validasi Manual)
	var product models.Product
	if err := config.DB.First(&product, txData.ProductID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
	}

	if float64(product.Stock) < txData.Quantity {
		return c.Status(400).JSON(fiber.Map{
			"error": "Stok tidak mencukupi",
			"stok_saat_ini": product.Stock,
			"permintaan": txData.Quantity,
		})
	}

	// Mulai Database Transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// A. Simpan data transaksi keluar
		if err := tx.Create(&txData).Error; err != nil {
			return err
		}

		// B. Update stok di tabel products (Stok - Quantity)
		if err := tx.Model(&models.Product{}).
			Where("id = ?", txData.ProductID).
			UpdateColumn("stock", gorm.Expr("stock - ?", txData.Quantity)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses transaksi keluar: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Barang keluar berhasil dicatat & stok berkurang",
		"data":    txData,
	})
}

// 4. Lihat Riwayat Barang Keluar
func GetTransactionsOut(c *fiber.Ctx) error {
	var transactions []models.TransactionOut
	config.DB.Preload("Product").Order("date desc").Find(&transactions)
	return c.JSON(transactions)
}