package handler

import (
	adminHandler "api-pizza-delivery/internal/handler/admin"
	authHandler "api-pizza-delivery/internal/handler/auth"
	user "api-pizza-delivery/internal/handler/user"
	"api-pizza-delivery/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	auth     *authHandler.AuthHandler
	product  *adminHandler.ProductHandler
	category *adminHandler.CategoryHandler
	admin    *adminHandler.AdminHandler
	upload   *adminHandler.UploadHandler
	profile  *user.ProfileHandler
	cart     *user.CartHandler
	order    *user.OrderHandler
	address  *user.AddressHandler
	favorite *user.FavoriteHandler
	review   *user.ReviewHandler
}

func NewRouter(
	auth *authHandler.AuthHandler,
	product *adminHandler.ProductHandler,
	category *adminHandler.CategoryHandler,
	admin *adminHandler.AdminHandler,
	upload *adminHandler.UploadHandler,
	profile *user.ProfileHandler,
	cart *user.CartHandler,
	order *user.OrderHandler,
	address *user.AddressHandler,
	favorite *user.FavoriteHandler,
	review *user.ReviewHandler,
) *Router {
	return &Router{
		auth:     auth,
		product:  product,
		category: category,
		admin:    admin,
		upload:   upload,
		profile:  profile,
		cart:     cart,
		order:    order,
		address:  address,
		favorite: favorite,
		review:   review,
	}
}

func (ro *Router) Setup(jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(corsMiddleware)
	// Глобальные middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Публичные маршруты
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", ro.auth.Register)
		r.Post("/login", ro.auth.Login)
		r.Post("/reset-request", ro.auth.ResetRequest)
		r.Post("/reset-password", ro.auth.ResetPassword)
	})

	// Публичные маршруты продуктов
	r.Get("/api/products", ro.product.GetAll)
	r.Get("/api/products/{id}", ro.product.GetByID)
	r.Get("/api/products/{id}/reviews", ro.review.GetByProductID)

	// Защищённые маршруты пользователя
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		// Auth
		r.Get("/api/auth/me", ro.auth.Me)
		r.Post("/api/auth/logout", ro.profile.Logout)

		// Профиль
		r.Get("/api/profile", ro.profile.GetProfile)
		r.Put("/api/profile", ro.profile.UpdateProfile)
		r.Put("/api/profile/password", ro.profile.ChangePassword)

		// Корзина
		r.Get("/api/cart", ro.cart.GetCart)
		r.Post("/api/cart", ro.cart.AddItem)
		r.Patch("/api/cart/{id}", ro.cart.UpdateItem)
		r.Delete("/api/cart/{id}", ro.cart.DeleteItem)

		// Заказы
		r.Get("/api/orders", ro.order.GetOrders)
		r.Get("/api/orders/{id}", ro.order.GetOrder)
		r.Post("/api/orders", ro.order.CreateOrder)

		// Адреса доставки
		r.Get("/api/delivery-addresses", ro.address.List)
		r.Post("/api/delivery-addresses", ro.address.Create)
		r.Put("/api/delivery-addresses/{id}", ro.address.Update)
		r.Delete("/api/delivery-addresses/{id}", ro.address.Delete)

		// Избранное
		r.Get("/api/favorites", ro.favorite.GetAll)
		r.Post("/api/favorites", ro.favorite.Add)
		r.Delete("/api/favorites/{id}", ro.favorite.Delete)

		// Отзывы
		r.Post("/api/products/{id}/reviews", ro.review.Create)
	})

	// Защищённые маршруты админа
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))
		r.Use(middleware.AdminOnly)

		// Пользователи
		r.Get("/api/admin/users", ro.admin.GetAllUsers)
		r.Get("/api/admin/users/export", ro.admin.ExportUsers)

		// Загрузка изображений (Cloudinary)
		r.Post("/api/admin/upload/image", ro.upload.UploadImage)

		// Пиццы
		r.Get("/api/admin/pizzas", ro.product.GetAllPizzas)
		r.Post("/api/admin/pizzas", ro.product.CreatePizza)
		r.Put("/api/admin/pizzas/{id}", ro.product.UpdatePizza)
		r.Delete("/api/admin/pizzas/{id}", ro.product.DeletePizza)

		// Напитки
		r.Get("/api/admin/drinks", ro.product.GetAllDrinks)
		r.Post("/api/admin/drinks", ro.product.CreateDrink)
		r.Put("/api/admin/drinks/{id}", ro.product.UpdateDrink)
		r.Delete("/api/admin/drinks/{id}", ro.product.DeleteDrink)

		// Категории товаров (допы)
		r.Get("/api/admin/product-categories", ro.category.List)
		r.Post("/api/admin/product-categories", ro.category.Create)
		r.Put("/api/admin/product-categories/{id}", ro.category.Update)
		r.Delete("/api/admin/product-categories/{id}", ro.category.Delete)

		// Допы
		r.Get("/api/admin/extras", ro.product.GetAllExtras)
		r.Post("/api/admin/extras", ro.product.CreateExtra)
		r.Put("/api/admin/extras/{id}", ro.product.UpdateExtra)
		r.Delete("/api/admin/extras/{id}", ro.product.DeleteExtra)

		// Отзывы
		r.Get("/api/admin/reviews", ro.admin.GetAllReviews)
		r.Delete("/api/admin/reviews/{id}", ro.admin.DeleteReview)

		// Аналитика
		r.Get("/api/admin/analytics", ro.admin.GetAnalytics)
		r.Get("/api/admin/analytics/export", ro.admin.ExportAnalytics)
	})

	return r
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight запрос
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}