package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"assignment_3/services"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", IndexHandler)

	return router
}

func IndexHandler(c *gin.Context) {
	data := services.GetStatusData()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Water":       data.Water,
		"Wind":        data.Wind,
		"WaterStatus": data.WaterStatus,
		"WindStatus":  data.WindStatus,
	})
}