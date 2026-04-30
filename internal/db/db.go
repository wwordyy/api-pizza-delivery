package db

import (
	"api-pizza-delivery/internal/models"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.ProductCategory{},
		&models.Role{},
		&models.User{},
		&models.PasswordResetCode{},
		&models.Product{},
		&models.SizePizza{},
		&models.ThicknessPizza{},
		&models.Pizza{},
		&models.TypeDrink{},
		&models.VolumeDrink{},
		&models.Drink{},
		&models.PaymentMethod{},
		&models.StatusOrder{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.DeliveryAddress{},
		&models.Favorite{},
		&models.Review{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := seedProductCategoriesAndMigrateExtras(db); err != nil {
		return fmt.Errorf("seed/migrate extras: %w", err)
	}

	return nil
}

func seedProductCategoriesAndMigrateExtras(db *gorm.DB) error {
	seed := []models.ProductCategory{
		{Title: "Снеки"},
		{Title: "Соусы"},
	}
	for i := range seed {
		var existing models.ProductCategory
		err := db.Where("title = ?", seed[i].Title).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&seed[i]).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	var snackCat, sauceCat models.ProductCategory
	if err := db.Where("title = ?", "Снеки").First(&snackCat).Error; err != nil {
		return err
	}
	if err := db.Where("title = ?", "Соусы").First(&sauceCat).Error; err != nil {
		return err
	}

	if err := db.Model(&models.Product{}).Where("product_type = ?", "snack").Updates(map[string]interface{}{
		"product_type": "extra",
		"category_id":  snackCat.ID,
	}).Error; err != nil {
		return err
	}
	if err := db.Model(&models.Product{}).Where("product_type = ?", "sauce").Updates(map[string]interface{}{
		"product_type": "extra",
		"category_id":  sauceCat.ID,
	}).Error; err != nil {
		return err
	}

	return nil
}

