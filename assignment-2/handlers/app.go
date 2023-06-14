package handlers

import (
	"assignment_2/database"
	"assignment_2/handlers/http_handlers"
	_ "assignment_2/docs"
	"assignment_2/repository/items_repository/item_pg"
	"assignment_2/repository/orders_repository/orders_pg"
	"assignment_2/services"
	
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"
)


// @title Order Items Documentation
// @version 1.0
// @description This is a server for Order Items.
// @termsOfService http://swagger.io/terms/
// @contact.name Swagger API Team
// @host localhost:3000
// @BasePath /
func StartApp() {
	database.InitiliazeDatabase()
	router := gin.Default()

	database.InitiliazeDatabase()

	db := database.GetDatabaseInstance()

	itemRepo := item_pg.NewItemPG(db)
	orderRepo := orders_pg.NewOrderPG(db)

	itemService := services.NewItemService(itemRepo)
	orderService := services.NewOrderService(orderRepo, itemService)

	orderHandler := http_handlers.NewOrderHandler(orderService)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := router.Group("/orders")
	{
		v1.POST("/", orderHandler.CreateOrder)
		v1.PUT("/:orderId", orderHandler.UpdateOrder)
		v1.GET("/", orderHandler.GetAllOrders)
		v1.GET("/:orderId", orderHandler.GetOrderById)
		v1.DELETE("/:orderId", orderHandler.DeleteOrders)
	}
	router.Run(":3000")
}

