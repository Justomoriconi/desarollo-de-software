package controllers

import (
	"backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetEventos(c *gin.Context) {
	eventos, error := services.GetEventos()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener eventos",
		})
		return
	}

	c.JSON(http.StatusOK, eventos)
}
