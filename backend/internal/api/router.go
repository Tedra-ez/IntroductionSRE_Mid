package api

import (
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/handlers"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/middleware"
	"github.com/Tedra-ez/AdvancedProgramming_Final/internal/services"
	"github.com/gin-gonic/gin"
)

func SetUpRouters(r *gin.Engine, orderHandler *handlers.OrderHandler, productHandler *handlers.ProductHandler, authHandler *handlers.AuthHandler, pageHandler *handlers.PageHandler, analyticsHandler *handlers.AnalyticsHandler, authSvc *services.AuthService) {
	r.Use(middleware.Logger(), middleware.CORS(), middleware.Auth(authSvc))

	r.GET("/", pageHandler.Index)
	r.GET("/shop", pageHandler.Shop)
	r.GET("/product/:id", pageHandler.Product)
	r.GET("/account", pageHandler.Account)
	r.GET("/wishlist", pageHandler.Wishlist)
	r.GET("/cart", pageHandler.Cart)
	r.GET("/checkout", pageHandler.Checkout)
	r.GET("/login", pageHandler.LoginPage)
	r.GET("/register", pageHandler.RegisterPage)
	r.GET("/account/orders", middleware.RequireAuth, pageHandler.AccountOrders)

	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth, middleware.RequireAdmin)
	{
		admin.GET("", pageHandler.AdminDashboard)
		admin.GET("/dashboard", pageHandler.AdminDashboard)
		admin.GET("/orders", pageHandler.AdminOrders)
		admin.GET("/products", pageHandler.AdminProducts)
		admin.GET("/users", pageHandler.AdminUsers)
		admin.GET("/users/:userId/orders", pageHandler.AdminUserOrders)
		admin.GET("/analytics", pageHandler.AdminAnalytics)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout", authHandler.Logout)
	}

	orders := r.Group("/orders")
	{
		orders.GET("", orderHandler.ListOrdersByUser)
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("/:id", orderHandler.GetOrderStatus)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	api := r.Group("/api")
	{
		api.GET("/product", productHandler.GetProducts)
		api.GET("/product/:id", productHandler.GetProductByID)

		analytics := api.Group("/analytics")
		analytics.Use(middleware.RequireAuth, middleware.RequireAdmin)
		{
			analytics.GET("/stats", analyticsHandler.DashboardStatsHandler())
			analytics.GET("/top-products", analyticsHandler.TopProductsHandler())
			analytics.GET("/revenue", analyticsHandler.RevenueHandler())
			analytics.GET("/orders-status", analyticsHandler.OrdersByStatusHandler())
		}

		adminAPI := api.Group("")
		adminAPI.Use(middleware.RequireAuth, middleware.RequireAdmin)
		adminAPI.POST("/product", productHandler.CreateProduct)
		adminAPI.PUT("/product/:id", productHandler.UpdateProduct)
		adminAPI.DELETE("/product/:id", productHandler.DeleteProduct)
	}
}
