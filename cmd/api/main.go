package main

import (
	"api-pizza-delivery/internal/handler"
	adminHr "api-pizza-delivery/internal/handler/admin"
	authHr "api-pizza-delivery/internal/handler/auth"
	userHr "api-pizza-delivery/internal/handler/user"
	adminRepo "api-pizza-delivery/internal/repository/admin"
	authRepo "api-pizza-delivery/internal/repository/auth"
	userRepo "api-pizza-delivery/internal/repository/user"
	adminService "api-pizza-delivery/internal/service/admin"
	authService "api-pizza-delivery/internal/service/auth"
	userService "api-pizza-delivery/internal/service/user"
	"fmt"
	"log"
	"net/http"

	"api-pizza-delivery/internal/config"
	"api-pizza-delivery/internal/db"
	"api-pizza-delivery/internal/mail"
	"api-pizza-delivery/internal/storage"
)

func main() {
	// 1. Конфиг
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// 2. Подключение к БД
	database, err := db.New(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	log.Println("connected to database")

	// 3. Миграции
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migration error: %v", err)
	}
	log.Println("migrations applied successfully")

	// 4. Репозитории
	authRepos    := authRepo.NewUserRepository(database)
	resetRepo   := authRepo.NewPasswordResetRepository(database)
	productRepo   := adminRepo.NewProductRepository(database)
	categoryRepo := adminRepo.NewCategoryRepository(database)

	adminUserRepo  := adminRepo.NewAdminUserRepository(database)
	adminReviewRepo := adminRepo.NewReviewRepository(database)
	analyticsRepo  := adminRepo.NewAnalyticsRepository(database)

	profileRepo  := userRepo.NewProfileRepository(database)
	cartRepo     := userRepo.NewCartRepository(database)
	orderRepo    := userRepo.NewOrderRepository(database)
	addressRepo  := userRepo.NewDeliveryAddressRepository(database)
	favoriteRepo := userRepo.NewFavoriteRepository(database)
	reviewRepo   := userRepo.NewReviewRepository(database)

	// 5. Сервисы
	mailCfg := mail.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		User:     cfg.SMTP.User,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}
	var mailSender mail.Sender
	if mailCfg.IsConfigured() {
		mailSender = mail.NewSMTP(mailCfg)
		log.Println("SMTP: включён")
	} else {
		log.Println("SMTP: выключен")
	}
	authSvc := authService.NewAuthService(authRepos, resetRepo, mailSender, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	productSvc  := adminService.NewProductService(productRepo, categoryRepo)
	categorySvc := adminService.NewCategoryService(categoryRepo)
	adminSvc    := adminService.NewAdminService(adminUserRepo, adminReviewRepo, analyticsRepo)

	profileSvc  := userService.NewProfileService(profileRepo)
	cartSvc     := userService.NewCartService(cartRepo)
	orderSvc    := userService.NewOrderService(orderRepo, cartRepo, profileRepo, mailSender, database)
	addressSvc  := userService.NewAddressService(addressRepo)
	favoriteSvc := userService.NewFavoriteService(favoriteRepo)
	reviewSvc   := userService.NewReviewService(reviewRepo)

	// 6. Хендлеры
	authHandler     := authHr.NewAuthHandler(authSvc)
	productHandler  := adminHr.NewProductHandler(productSvc)
	categoryHandler := adminHr.NewCategoryHandler(categorySvc)
	adminHandler    := adminHr.NewAdminHandler(adminSvc)

	var uploadHandler *adminHr.UploadHandler
	if cfg.Cloudinary.IsConfigured() {
		cld, err := storage.New(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
		if err != nil {
			log.Fatalf("cloudinary: %v", err)
		}
		uploadHandler = adminHr.NewUploadHandler(cld)
	} else {
		uploadHandler = adminHr.NewUploadHandler(nil)
	}

	profileHandler  := userHr.NewProfileHandler(profileSvc)
	cartHandler     := userHr.NewCartHandler(cartSvc)
	orderHandler    := userHr.NewOrderHandler(orderSvc)
	addressHandler  := userHr.NewAddressHandler(addressSvc)
	favoriteHandler := userHr.NewFavoriteHandler(favoriteSvc)
	reviewHandler   := userHr.NewReviewHandler(reviewSvc)


	router := handler.NewRouter(
		authHandler,
		productHandler,
		categoryHandler,
		adminHandler,
		uploadHandler,
		profileHandler,
		cartHandler,
		orderHandler,
		addressHandler,
		favoriteHandler,
		reviewHandler,
	)
	r := router.Setup(cfg.JWT.Secret)


	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}