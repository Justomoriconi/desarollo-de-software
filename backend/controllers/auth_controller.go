package controllers

import (
	"backend/dto"
	"backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var request dto.RegisterRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos invalidos",
		})

		return
	}

	id, err := services.Register(
		request.Nombre,
		request.Email,
		request.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al registrar  usuario",
		})
	}

	c.JSON(http.StatusCreated, id)

}

func Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos invalidos",
		})

		return
	}

	token, err := services.Login(
		request.Email,
		request.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
